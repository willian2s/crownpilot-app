using System.Net.Http.Headers;
using System.Security.Claims;
using System.Text.Encodings.Web;
using CrownPilot.Api.Configuration;
using CrownPilot.Application.Identity;
using Microsoft.AspNetCore.Authentication;
using Microsoft.Extensions.Options;
using Microsoft.Net.Http.Headers;

namespace CrownPilot.Api.Authentication;

public static class ContractAuthenticationDefaults
{
    public const string Scheme = "Bearer";
    public const string BootstrapPolicy = "BootstrapContractAccess";
    public const string PermissionClaim = "contract_permission";
    public const string BootstrapReadPermission = "bootstrap:read";

    // Deterministic local fixtures prove pipeline behavior without pretending to verify Firebase.
    public const string AuthenticatedToken = "contract-authenticated-token";
    public const string AuthorizedToken = "contract-authorized-token";
}

public sealed class ContractBearerAuthenticationHandler(
    IOptionsMonitor<AuthenticationSchemeOptions> options,
    ILoggerFactory logger,
    UrlEncoder encoder,
    IOptionsMonitor<RuntimeOptions> runtimeOptions,
    IFirebaseTokenVerifier firebaseTokenVerifier)
    : AuthenticationHandler<AuthenticationSchemeOptions>(options, logger, encoder)
{
    protected override async Task<AuthenticateResult> HandleAuthenticateAsync()
    {
        if (!Request.Headers.TryGetValue(HeaderNames.Authorization, out var authorizationValues) ||
            authorizationValues.Count != 1 ||
            !AuthenticationHeaderValue.TryParse(authorizationValues[0], out var header) ||
            !string.Equals(header.Scheme, ContractAuthenticationDefaults.Scheme,
                StringComparison.OrdinalIgnoreCase) ||
            string.IsNullOrWhiteSpace(header.Parameter))
        {
            return AuthenticateResult.NoResult();
        }

        var runtime = runtimeOptions.CurrentValue;
        if (runtime.UseContractAuthentication)
        {
            return AuthenticateContractFixture(header.Parameter);
        }

        if (!runtime.UseFirebaseAuthentication)
        {
            return AuthenticateResult.Fail("Bearer authentication is not configured.");
        }

        FirebaseVerifiedToken? verifiedToken;
        try
        {
            verifiedToken = await firebaseTokenVerifier.VerifyAsync(
                header.Parameter,
                Context.RequestAborted);
        }
        catch (OperationCanceledException) when (Context.RequestAborted.IsCancellationRequested)
        {
            throw;
        }
        catch (Exception)
        {
            return AuthenticateResult.Fail("Bearer token is invalid.");
        }

        if (verifiedToken is null || string.IsNullOrWhiteSpace(verifiedToken.Uid))
        {
            return AuthenticateResult.Fail("Bearer token is invalid.");
        }

        return AuthenticateResult.Success(CreateTicket(verifiedToken.Uid, false));
    }

    private static AuthenticateResult AuthenticateContractFixture(string token)
    {
        var hasAuthentication = token == ContractAuthenticationDefaults.AuthenticatedToken;
        var hasAuthorization = token == ContractAuthenticationDefaults.AuthorizedToken;

        if (!hasAuthentication && !hasAuthorization)
        {
            return AuthenticateResult.Fail("Bearer token is invalid.");
        }

        var claims = new Dictionary<string, object>();
        if (hasAuthorization)
        {
            claims[ContractAuthenticationDefaults.PermissionClaim] =
                ContractAuthenticationDefaults.BootstrapReadPermission;
        }

        return AuthenticateResult.Success(CreateTicket("contract-firebase-uid", hasAuthorization));
    }

    private static AuthenticationTicket CreateTicket(
        string firebaseUid,
        bool hasAuthorization)
    {
        var claims = new List<Claim>
        {
            // These claims are derived from the verifier result, never request input.
            new("sub", firebaseUid),
            new(ClaimTypes.NameIdentifier, firebaseUid),
            new("firebase_uid", firebaseUid)
        };

        if (hasAuthorization)
        {
            claims.Add(new Claim(
                ContractAuthenticationDefaults.PermissionClaim,
                ContractAuthenticationDefaults.BootstrapReadPermission));
        }

        var identity = new ClaimsIdentity(claims, ContractAuthenticationDefaults.Scheme);
        return new AuthenticationTicket(
            new ClaimsPrincipal(identity),
            ContractAuthenticationDefaults.Scheme);
    }

    protected override Task HandleChallengeAsync(AuthenticationProperties properties)
    {
        Response.StatusCode = StatusCodes.Status401Unauthorized;
        Response.Headers.WWWAuthenticate = ContractAuthenticationDefaults.Scheme;
        return Task.CompletedTask;
    }

    protected override Task HandleForbiddenAsync(AuthenticationProperties properties)
    {
        Response.StatusCode = StatusCodes.Status403Forbidden;
        return Task.CompletedTask;
    }
}
