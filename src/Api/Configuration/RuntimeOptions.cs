using Microsoft.AspNetCore.Cors.Infrastructure;

namespace CrownPilot.Api.Configuration;

public sealed class RuntimeOptions
{
    public const string SectionName = "CrownPilot";

    public string? Environment { get; set; }

    public CorsOptions Cors { get; set; } = new();

    public DocumentationOptions Documentation { get; set; } = new();

    public AuthenticationBoundaryOptions Authentication { get; set; } = new();

    public RuntimeEnvironment ResolvedEnvironment =>
        RuntimeEnvironmentNames.TryParse(Environment, out var environment)
            ? environment
            : throw new InvalidOperationException("CrownPilot environment has not been resolved.");

    public bool ExposeOpenApiJson =>
        (ResolvedEnvironment is RuntimeEnvironment.Local or RuntimeEnvironment.Staging) &&
        Documentation.ExposeJson is true;

    public bool ExposeOpenApiUi =>
        ExposeOpenApiJson && Documentation.ExposeUi is true;

    public bool UseContractAuthentication =>
        ResolvedEnvironment is RuntimeEnvironment.Local &&
        Authentication.ContractFixturesEnabled is true &&
        string.Equals(Authentication.Mode, AuthenticationBoundaryOptions.ContractFixtureMode,
            StringComparison.OrdinalIgnoreCase);

    public bool UseFirebaseAuthentication =>
        string.Equals(Authentication.Mode, AuthenticationBoundaryOptions.FirebaseMode,
            StringComparison.OrdinalIgnoreCase) &&
        Authentication.Firebase.IsConfigured;

    public void ResolveForHost(string hostEnvironment)
    {
        if (!RuntimeEnvironmentNames.TryParseHost(hostEnvironment, out var hostRuntimeEnvironment))
        {
            throw new InvalidOperationException(
                $"Unsupported ASP.NET Core environment '{hostEnvironment}'.");
        }

        if (!string.IsNullOrWhiteSpace(Environment))
        {
            if (!RuntimeEnvironmentNames.TryParse(Environment, out var configuredRuntimeEnvironment) ||
                configuredRuntimeEnvironment != hostRuntimeEnvironment)
            {
                throw new InvalidOperationException(
                    "CrownPilot:Environment must match ASPNETCORE_ENVIRONMENT.");
            }
        }

        Environment = hostRuntimeEnvironment.ToString();
        Cors ??= new CorsOptions();
        Documentation ??= new DocumentationOptions();
        Authentication ??= new AuthenticationBoundaryOptions();
        Authentication.Firebase ??= new FirebaseAuthenticationBoundaryOptions();

        Cors.AllowedOrigins ??= hostRuntimeEnvironment is RuntimeEnvironment.Local
            ? ["http://localhost:5173", "http://127.0.0.1:5173"]
            : [];

        Documentation.ExposeJson ??= hostRuntimeEnvironment is RuntimeEnvironment.Local or RuntimeEnvironment.Staging;
        Documentation.ExposeUi ??= hostRuntimeEnvironment is RuntimeEnvironment.Local or RuntimeEnvironment.Staging;

        Authentication.Mode ??= hostRuntimeEnvironment is RuntimeEnvironment.Local
            ? AuthenticationBoundaryOptions.ContractFixtureMode
            : AuthenticationBoundaryOptions.UnconfiguredMode;
        Authentication.ContractFixturesEnabled ??= hostRuntimeEnvironment is RuntimeEnvironment.Local;
    }

    public CorsPolicy BuildCorsPolicy()
    {
        var builder = new CorsPolicyBuilder();
        ConfigureCorsPolicy(builder);
        return builder.Build();
    }

    public void ConfigureCorsPolicy(CorsPolicyBuilder builder)
    {
        var origins = Cors.AllowedOrigins ?? [];

        if (origins.Length == 0)
        {
            builder.SetIsOriginAllowed(_ => false);
        }
        else
        {
            builder.WithOrigins(origins);
        }

        builder
            .WithMethods("GET", "HEAD", "OPTIONS", "PUT", "DELETE")
            .WithHeaders("Accept", "Authorization", "Content-Type")
            ;
    }
}

public sealed class CorsOptions
{
    public string[]? AllowedOrigins { get; set; }
}

public sealed class DocumentationOptions
{
    public bool? ExposeJson { get; set; }

    public bool? ExposeUi { get; set; }
}

public sealed class AuthenticationBoundaryOptions
{
    public const string ContractFixtureMode = "ContractFixture";
    public const string FirebaseMode = "Firebase";
    public const string UnconfiguredMode = "Unconfigured";

    public string? Mode { get; set; }

    public bool? ContractFixturesEnabled { get; set; }

    public FirebaseAuthenticationBoundaryOptions Firebase { get; set; } = new();
}

public sealed class FirebaseAuthenticationBoundaryOptions
{
    public string? ProjectId { get; set; }

    public string? Issuer { get; set; }

    // Server-only runtime secret sources. ADC remains preferred when both are empty.
    public string? ServiceAccountJson { get; set; }

    public string? ServiceAccountFile { get; set; }

    public bool IsConfigured =>
        !string.IsNullOrWhiteSpace(ProjectId) &&
        string.Equals(Issuer, ExpectedIssuer, StringComparison.Ordinal);

    public bool HasAnyConfiguration =>
        !string.IsNullOrWhiteSpace(ProjectId) ||
        !string.IsNullOrWhiteSpace(Issuer) ||
        !string.IsNullOrWhiteSpace(ServiceAccountJson) ||
        !string.IsNullOrWhiteSpace(ServiceAccountFile);

    public string ExpectedIssuer =>
        string.IsNullOrWhiteSpace(ProjectId)
            ? string.Empty
            : $"https://securetoken.google.com/{ProjectId}";
}
