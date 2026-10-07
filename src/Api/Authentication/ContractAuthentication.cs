using System.Net.Http.Headers;
using System.Security.Claims;
using System.Text.Encodings.Web;
using CrownPilot.Api.Configuration;
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
    IOptionsMonitor<RuntimeOptions> runtimeOptions)
    : AuthenticationHandler<AuthenticationSchemeOptions>(options, logger, encoder)
{
    protected override Task<AuthenticateResult> HandleAuthenticateAsync()
    {
        if (!Request.Headers.TryGetValue(HeaderNames.Authorization, out var authorizationValues) ||
            authorizationValues.Count != 1 ||
            !AuthenticationHeaderValue.TryParse(authorizationValues[0], out var header) ||
            !string.Equals(header.Scheme, ContractAuthenticationDefaults.Scheme,
                StringComparison.OrdinalIgnoreCase) ||
            string.IsNullOrWhiteSpace(header.Parameter))
        {
            return Task.FromResult(AuthenticateResult.NoResult());
        }

        var runtime = runtimeOptions.CurrentValue;
        if (!runtime.UseContractAuthentication)
        {
            return Task.FromResult(AuthenticateResult.Fail("Bearer provider is not configured."));
        }

        var token = header.Parameter;
        var hasAuthentication = token == ContractAuthenticationDefaults.AuthenticatedToken;
        var hasAuthorization = token == ContractAuthenticationDefaults.AuthorizedToken;

        if (!hasAuthentication && !hasAuthorization)
        {
            return Task.FromResult(AuthenticateResult.Fail("Bearer token is invalid."));
        }

        var claims = new List<Claim>
        {
            new("sub", "contract-firebase-uid"),
            new(ClaimTypes.NameIdentifier, "contract-firebase-uid"),
            new("firebase_uid", "contract-firebase-uid")
        };

        if (hasAuthorization)
        {
            claims.Add(new Claim(
                ContractAuthenticationDefaults.PermissionClaim,
                ContractAuthenticationDefaults.BootstrapReadPermission));
        }

        var identity = new ClaimsIdentity(claims, ContractAuthenticationDefaults.Scheme);
        return Task.FromResult(AuthenticateResult.Success(new AuthenticationTicket(
            new ClaimsPrincipal(identity),
            ContractAuthenticationDefaults.Scheme)));
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
