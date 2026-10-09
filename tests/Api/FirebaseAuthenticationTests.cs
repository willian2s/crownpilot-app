using System.Net;
using System.Net.Http.Headers;
using System.Text.Json;
using CrownPilot.Api.Authentication;
using CrownPilot.Application.Identity;
using CrownPilot.Infrastructure.Authentication;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Mvc.Testing;
using Microsoft.AspNetCore.TestHost;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging.Abstractions;
using Microsoft.Extensions.Options;

namespace Api.Tests;

[Trait("Category", "Authentication")]
public sealed class FirebaseAuthenticationTests
{
    [Fact]
    public async Task MissingBearerReturnsGeneric401()
    {
        await using var factory = CreateFactory(new FakeFirebaseTokenVerifier());
        using var client = factory.CreateClient();

        using var response = await client.GetAsync("/api/v1/bootstrap");
        var body = await response.Content.ReadAsStringAsync();

        Assert.Equal(HttpStatusCode.Unauthorized, response.StatusCode);
        Assert.Equal("Bearer", response.Headers.WwwAuthenticate.Single().Scheme);
        Assert.Contains("authentication_required", body, StringComparison.Ordinal);
        Assert.DoesNotContain("Firebase", body, StringComparison.OrdinalIgnoreCase);
        Assert.DoesNotContain("UID", body, StringComparison.OrdinalIgnoreCase);
    }

    [Theory]
    [InlineData("invalid")]
    [InlineData("expired")]
    [InlineData("wrong-issuer")]
    [InlineData("wrong-audience")]
    [InlineData("wrong-project")]
    [InlineData("bad-signature")]
    [InlineData("rotated-kid-rejected")]
    [InlineData("empty-sub")]
    public async Task InvalidFirebaseTokensReturnSame401(string token)
    {
        var verifier = new FakeFirebaseTokenVerifier();
        await using var factory = CreateFactory(verifier);
        using var client = factory.CreateClient();
        client.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Bearer", token);

        using var response = await client.GetAsync("/api/v1/bootstrap");
        var body = await response.Content.ReadAsStringAsync();

        Assert.Equal(HttpStatusCode.Unauthorized, response.StatusCode);
        Assert.Contains("authentication_required", body, StringComparison.Ordinal);
        Assert.DoesNotContain(token, body, StringComparison.Ordinal);
    }

    [Fact]
    [Trait("Category", "Authorization")]
    public async Task VerifiedFirebaseTokenCanAuthenticateWithoutGrantingAuthorization()
    {
        var verifier = new FakeFirebaseTokenVerifier();
        verifier.Add("authenticated-only", new FirebaseVerifiedToken(
            "verified-uid",
            "https://securetoken.google.com/test-project",
            "test-project"));
        await using var factory = CreateFactory(verifier);
        using var client = factory.CreateClient();
        client.DefaultRequestHeaders.Authorization =
            new AuthenticationHeaderValue("Bearer", "authenticated-only");

        using var response = await client.GetAsync("/api/v1/bootstrap");
        var problem = await JsonSerializer.DeserializeAsync<JsonElement>(
            await response.Content.ReadAsStreamAsync());

        Assert.Equal(HttpStatusCode.Forbidden, response.StatusCode);
        Assert.Equal("authorization_forbidden", problem.GetProperty("code").GetString());
    }

    [Fact]
    public async Task VerifiedUidCannotBeReplacedByRequestFields()
    {
        var verifier = new FakeFirebaseTokenVerifier();
        verifier.Add("valid", new FirebaseVerifiedToken(
            "verified-uid",
            "https://securetoken.google.com/test-project",
            "test-project"));
        await using var factory = CreateFactory(verifier);
        using var client = factory.CreateClient();

        using var request = new HttpRequestMessage(
            HttpMethod.Get,
            "/api/v1/bootstrap?firebase_uid=attacker-uid");
        request.Headers.Authorization = new AuthenticationHeaderValue("Bearer", "valid");
        request.Headers.Add("X-Firebase-Uid", "attacker-uid");
        request.Content = new StringContent("{\"firebaseUid\":\"attacker-uid\"}");

        using var response = await client.SendAsync(request);
        var body = await response.Content.ReadAsStringAsync();

        Assert.Equal(HttpStatusCode.Forbidden, response.StatusCode);
        Assert.Contains("authorization_forbidden", body, StringComparison.Ordinal);
        Assert.DoesNotContain("attacker-uid", body, StringComparison.Ordinal);
        Assert.Equal("valid", verifier.LastToken);
    }

    [Fact]
    public async Task ContractTokenIsNotAcceptedOutsideLocalFixtureMode()
    {
        var verifier = new FakeFirebaseTokenVerifier();
        await using var factory = CreateFactory(verifier, "Staging");
        using var client = factory.CreateClient();
        client.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue(
            "Bearer",
            ContractAuthenticationDefaults.AuthorizedToken);

        using var response = await client.GetAsync("/api/v1/bootstrap");

        Assert.Equal(HttpStatusCode.Unauthorized, response.StatusCode);
        Assert.Equal(ContractAuthenticationDefaults.AuthorizedToken, verifier.LastToken);
    }

    [Fact]
    public void FirebaseConfigurationBindsProjectAndExpectedIssuer()
    {
        var options = new FirebaseAuthenticationOptions
        {
            ProjectId = "test-project",
            Issuer = "https://securetoken.google.com/test-project"
        };

        Assert.Empty(FirebaseAuthenticationOptions.GetValidationErrors(options));
        Assert.True(options.IsConfigured);
    }

    [Fact]
    public void FirebaseConfigurationRejectsIssuerAndCredentialAmbiguity()
    {
        var options = new FirebaseAuthenticationOptions
        {
            ProjectId = "test-project",
            Issuer = "https://securetoken.google.com/other-project",
            ServiceAccountJson = "runtime-secret",
            ServiceAccountFile = "/runtime/secret.json"
        };

        var errors = FirebaseAuthenticationOptions.GetValidationErrors(options);

        Assert.Equal(2, errors.Count);
        Assert.All(errors, error => Assert.DoesNotContain("runtime-secret", error,
            StringComparison.Ordinal));
    }

    [Fact]
    public void RenderAdminCredentialFieldsMustBeComplete()
    {
        var options = new FirebaseAuthenticationOptions
        {
            AdminProjectId = "test-project",
            AdminPrivateKey = "runtime-secret"
        };

        var errors = FirebaseAuthenticationOptions.GetValidationErrors(options);

        Assert.Single(errors);
        Assert.DoesNotContain("runtime-secret", errors[0], StringComparison.Ordinal);
    }

    [Fact]
    public async Task OfficialVerifierFailsClosedBeforeProviderInitializationWhenUnconfigured()
    {
        var verifier = new FirebaseTokenVerifier(
            Options.Create(new FirebaseAuthenticationOptions()),
            NullLogger<FirebaseTokenVerifier>.Instance);

        Assert.Null(await verifier.VerifyAsync("fixture-token"));
        verifier.Dispose();
    }

    private static WebApplicationFactory<Program> CreateFactory(
        FakeFirebaseTokenVerifier verifier,
        string environment = "Development") =>
        new WebApplicationFactory<Program>()
            .WithWebHostBuilder(builder =>
            {
                builder.UseSetting(WebHostDefaults.EnvironmentKey, environment);
                builder.ConfigureAppConfiguration((_, configuration) =>
                    configuration.AddInMemoryCollection(new Dictionary<string, string?>
                    {
                        ["CrownPilot:Authentication:Mode"] = "Firebase",
                        ["CrownPilot:Authentication:ContractFixturesEnabled"] = "false",
                        ["CrownPilot:Authentication:Firebase:ProjectId"] = "test-project",
                        ["CrownPilot:Authentication:Firebase:Issuer"] =
                            "https://securetoken.google.com/test-project"
                    }));
                builder.ConfigureTestServices(services =>
                    services.AddSingleton<IFirebaseTokenVerifier>(verifier));
            });

    private sealed class FakeFirebaseTokenVerifier : IFirebaseTokenVerifier
    {
        private readonly Dictionary<string, FirebaseVerifiedToken> tokens = new();

        public string? LastToken { get; private set; }

        public void Add(string token, FirebaseVerifiedToken verifiedToken) =>
            tokens[token] = verifiedToken;

        public Task<FirebaseVerifiedToken?> VerifyAsync(
            string idToken,
            CancellationToken cancellationToken = default)
        {
            LastToken = idToken;
            tokens.TryGetValue(idToken, out var verifiedToken);
            return Task.FromResult(verifiedToken);
        }
    }
}
