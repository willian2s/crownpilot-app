using Microsoft.Extensions.Options;

namespace CrownPilot.Api.Configuration;

public sealed class RuntimeOptionsValidator : IValidateOptions<RuntimeOptions>
{
    public ValidateOptionsResult Validate(string? name, RuntimeOptions options)
    {
        var errors = GetValidationErrors(options);
        return errors.Count == 0
            ? ValidateOptionsResult.Success
            : ValidateOptionsResult.Fail(errors);
    }

    public static IReadOnlyList<string> GetValidationErrors(RuntimeOptions options)
    {
        var errors = new List<string>();

        if (!RuntimeEnvironmentNames.TryParse(options.Environment, out var environment))
        {
            errors.Add("CrownPilot:Environment must be Local, Preview, Staging, or Production.");
            return errors;
        }

        var origins = options.Cors?.AllowedOrigins ?? [];
        var duplicateOrigins = origins
            .GroupBy(origin => origin, StringComparer.OrdinalIgnoreCase)
            .Where(group => group.Count() > 1)
            .Select(group => group.Key)
            .ToArray();

        if (duplicateOrigins.Length > 0)
        {
            errors.Add("CrownPilot:Cors:AllowedOrigins cannot contain duplicates.");
        }

        foreach (var origin in origins)
        {
            if (!IsAllowedOrigin(origin, environment))
            {
                errors.Add($"CrownPilot:Cors:AllowedOrigins contains invalid origin '{origin}'.");
            }
        }

        var documentation = options.Documentation ?? new DocumentationOptions();
        if (documentation.ExposeUi is true && documentation.ExposeJson is not true)
        {
            errors.Add("CrownPilot:Documentation:ExposeUi requires ExposeJson.");
        }

        if (environment is RuntimeEnvironment.Preview or RuntimeEnvironment.Production &&
            (documentation.ExposeJson is true || documentation.ExposeUi is true))
        {
            errors.Add("OpenAPI exposure is disabled by policy in Preview and Production.");
        }

        var authentication = options.Authentication ?? new AuthenticationBoundaryOptions();
        if (authentication.ContractFixturesEnabled is true &&
            (environment is not RuntimeEnvironment.Local ||
             !string.Equals(authentication.Mode, AuthenticationBoundaryOptions.ContractFixtureMode,
                 StringComparison.OrdinalIgnoreCase)))
        {
            errors.Add("Contract authentication fixtures are allowed only in Local.");
        }

        var firebase = authentication.Firebase ?? new FirebaseAuthenticationBoundaryOptions();
        if (string.Equals(authentication.Mode, AuthenticationBoundaryOptions.FirebaseMode,
                StringComparison.OrdinalIgnoreCase))
        {
            var missingProject = string.IsNullOrWhiteSpace(firebase.ProjectId);
            var missingIssuer = string.IsNullOrWhiteSpace(firebase.Issuer);
            if (missingProject || missingIssuer)
            {
                errors.Add(missingProject && missingIssuer
                    ? "Firebase project ID and issuer are required when Firebase authentication is enabled."
                    : "Firebase project ID and issuer must be configured together.");
            }
            else if (!string.Equals(firebase.Issuer, firebase.ExpectedIssuer,
                         StringComparison.Ordinal))
            {
                errors.Add("Firebase issuer must match the configured project ID.");
            }

            if (!string.IsNullOrWhiteSpace(firebase.ServiceAccountJson) &&
                !string.IsNullOrWhiteSpace(firebase.ServiceAccountFile))
            {
                errors.Add("Firebase service account JSON and file cannot both be configured.");
            }
        }

        return errors;
    }

    private static bool IsAllowedOrigin(string? value, RuntimeEnvironment environment)
    {
        if (string.IsNullOrWhiteSpace(value) || value != value.Trim() || value == "*")
        {
            return false;
        }

        if (!Uri.TryCreate(value, UriKind.Absolute, out var uri) ||
            string.IsNullOrEmpty(uri.Host) ||
            uri.UserInfo.Length > 0 ||
            uri.AbsolutePath != "/" ||
            !string.IsNullOrEmpty(uri.Query) ||
            !string.IsNullOrEmpty(uri.Fragment) ||
            !string.Equals(uri.Scheme, Uri.UriSchemeHttp, StringComparison.OrdinalIgnoreCase) &&
            !string.Equals(uri.Scheme, Uri.UriSchemeHttps, StringComparison.OrdinalIgnoreCase))
        {
            return false;
        }

        if (environment is not RuntimeEnvironment.Local &&
            !string.Equals(uri.Scheme, Uri.UriSchemeHttps, StringComparison.OrdinalIgnoreCase))
        {
            return false;
        }

        if (string.Equals(uri.Scheme, Uri.UriSchemeHttp, StringComparison.OrdinalIgnoreCase) &&
            !IsLocalHost(uri.Host))
        {
            return false;
        }

        return true;
    }

    private static bool IsLocalHost(string host) =>
        host.Equals("localhost", StringComparison.OrdinalIgnoreCase) ||
        host.Equals("127.0.0.1", StringComparison.OrdinalIgnoreCase) ||
        host.Equals("[::1]", StringComparison.OrdinalIgnoreCase) ||
        host.Equals("::1", StringComparison.OrdinalIgnoreCase);
}
