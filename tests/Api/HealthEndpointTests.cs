using System.Net;
using System.Net.Http.Headers;
using System.Net.Http.Json;
using System.Text.Json;
using CrownPilot.Api.Authentication;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Mvc.Testing;

namespace Api.Tests;

public sealed class HealthEndpointTests(WebApplicationFactory<Program> factory)
    : IClassFixture<WebApplicationFactory<Program>>
{
    private readonly HttpClient client = factory.CreateClient();

    [Fact]
    public async Task LiveHealthDoesNotRequireExternalProvider()
    {
        using var response = await client.GetAsync("/health/live");

        Assert.Equal(HttpStatusCode.OK, response.StatusCode);
        Assert.Equal("Healthy", await response.Content.ReadAsStringAsync());
    }

    [Fact]
    public async Task OpenApiIsGeneratedByApiCode()
    {
        using var response = await client.GetAsync("/openapi/v1.json");
        using var document = await JsonDocument.ParseAsync(await response.Content.ReadAsStreamAsync());

        Assert.Equal(HttpStatusCode.OK, response.StatusCode);
        Assert.StartsWith("3.", document.RootElement.GetProperty("openapi").GetString());
        Assert.True(document.RootElement.GetProperty("paths").TryGetProperty("/api/v1/bootstrap", out _));
    }

    [Fact]
    public async Task SwaggerUiServesGeneratedContractInDevelopment()
    {
        using var response = await client.GetAsync("/docs/index.html");
        var body = await response.Content.ReadAsStringAsync();

        Assert.Equal(HttpStatusCode.OK, response.StatusCode);
        Assert.Contains("id=\"swagger-ui\"", body, StringComparison.Ordinal);
        Assert.Contains("CrownPilot API documentation", body, StringComparison.Ordinal);
    }

    [Fact]
    public async Task SwaggerUiServesGeneratedContractInStaging()
    {
        await using var factory = new WebApplicationFactory<Program>()
            .WithWebHostBuilder(builder =>
                builder.UseSetting(WebHostDefaults.EnvironmentKey, "Staging"));
        using var stagingClient = factory.CreateClient();

        using var response = await stagingClient.GetAsync("/docs/index.html");
        using var openApiResponse = await stagingClient.GetAsync("/openapi/v1.json");

        Assert.Equal(HttpStatusCode.OK, response.StatusCode);
        Assert.Equal(HttpStatusCode.OK, openApiResponse.StatusCode);
    }

    [Fact]
    public async Task PreviewDoesNotExposeDocumentationRoutes()
    {
        await using var factory = new WebApplicationFactory<Program>()
            .WithWebHostBuilder(builder =>
                builder.UseSetting(WebHostDefaults.EnvironmentKey, "Preview"));
        using var previewClient = factory.CreateClient();

        using var openApiResponse = await previewClient.GetAsync("/openapi/v1.json");
        using var docsResponse = await previewClient.GetAsync("/docs/index.html");

        Assert.Equal(HttpStatusCode.NotFound, openApiResponse.StatusCode);
        Assert.Equal(HttpStatusCode.NotFound, docsResponse.StatusCode);
    }

    [Fact]
    public async Task BootstrapRequiresAuthenticationAndAuthorization()
    {
        using var unauthenticatedResponse = await client.GetAsync("/api/v1/bootstrap");
        Assert.Equal(HttpStatusCode.Unauthorized, unauthenticatedResponse.StatusCode);
        Assert.Equal("Bearer", unauthenticatedResponse.Headers.WwwAuthenticate.Single().Scheme);
        var unauthorizedProblem = await unauthenticatedResponse.Content.ReadFromJsonAsync<JsonElement>();
        Assert.Equal("authentication_required", unauthorizedProblem.GetProperty("code").GetString());

        client.DefaultRequestHeaders.Authorization =
            new AuthenticationHeaderValue("Bearer", ContractAuthenticationDefaults.AuthenticatedToken);

        using var forbiddenResponse = await client.GetAsync("/api/v1/bootstrap");
        Assert.Equal(HttpStatusCode.Forbidden, forbiddenResponse.StatusCode);
        var forbiddenProblem = await forbiddenResponse.Content.ReadFromJsonAsync<JsonElement>();
        Assert.Equal("authorization_forbidden", forbiddenProblem.GetProperty("code").GetString());

        client.DefaultRequestHeaders.Authorization =
            new AuthenticationHeaderValue("Bearer", ContractAuthenticationDefaults.AuthorizedToken);
        using var authorizedResponse = await client.GetAsync("/api/v1/bootstrap");
        Assert.Equal(HttpStatusCode.OK, authorizedResponse.StatusCode);
    }

    [Fact]
    public async Task CorsAllowsOnlyConfiguredLocalOrigins()
    {
        using var allowedRequest = new HttpRequestMessage(HttpMethod.Get, "/health/live");
        allowedRequest.Headers.TryAddWithoutValidation("Origin", "http://localhost:5173");
        using var allowedResponse = await client.SendAsync(allowedRequest);

        Assert.Equal(HttpStatusCode.OK, allowedResponse.StatusCode);
        Assert.Equal("http://localhost:5173",
            allowedResponse.Headers.GetValues("Access-Control-Allow-Origin").Single());

        using var deniedRequest = new HttpRequestMessage(HttpMethod.Get, "/health/live");
        deniedRequest.Headers.TryAddWithoutValidation("Origin", "https://untrusted.example");
        using var deniedResponse = await client.SendAsync(deniedRequest);

        Assert.Equal(HttpStatusCode.OK, deniedResponse.StatusCode);
        Assert.False(deniedResponse.Headers.Contains("Access-Control-Allow-Origin"));

        using var preflightRequest = new HttpRequestMessage(HttpMethod.Options, "/health/live");
        preflightRequest.Headers.TryAddWithoutValidation("Origin", "http://localhost:5173");
        preflightRequest.Headers.TryAddWithoutValidation("Access-Control-Request-Method", "GET");
        preflightRequest.Headers.TryAddWithoutValidation("Access-Control-Request-Headers", "Authorization");
        using var preflightResponse = await client.SendAsync(preflightRequest);

        Assert.Equal(HttpStatusCode.NoContent, preflightResponse.StatusCode);
        Assert.Equal("http://localhost:5173",
            preflightResponse.Headers.GetValues("Access-Control-Allow-Origin").Single());
    }

    [Fact]
    public async Task ProductionDoesNotExposeDocumentationRoutes()
    {
        await using var factory = new WebApplicationFactory<Program>()
            .WithWebHostBuilder(builder =>
                builder.UseSetting(WebHostDefaults.EnvironmentKey, "Production"));
        using var productionClient = factory.CreateClient();

        using var openApiResponse = await productionClient.GetAsync("/openapi/v1.json");
        using var docsResponse = await productionClient.GetAsync("/docs/index.html");

        Assert.Equal(HttpStatusCode.NotFound, openApiResponse.StatusCode);
        Assert.Equal(HttpStatusCode.NotFound, docsResponse.StatusCode);
    }

    [Fact]
    public async Task MissingRouteUsesProblemDetailsWithoutInternalDetails()
    {
        using var response = await client.GetAsync("/not-found");
        var problem = await response.Content.ReadFromJsonAsync<JsonElement>();

        Assert.Equal(HttpStatusCode.NotFound, response.StatusCode);
        Assert.Equal("resource_not_found", problem.GetProperty("code").GetString());
        Assert.True(problem.TryGetProperty("traceId", out _));
        Assert.DoesNotContain("ConnectionStrings", problem.ToString(), StringComparison.OrdinalIgnoreCase);
    }
}
