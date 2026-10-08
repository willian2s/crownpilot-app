namespace CrownPilot.Infrastructure.Authentication;

public sealed class FirebaseAuthenticationOptions
{
    public const string SectionName = "CrownPilot:Authentication:Firebase";

    public string? ProjectId { get; set; }

    public string? Issuer { get; set; }

    // These values are supplied only by runtime secret storage. ADC is preferred
    // when neither explicit credential source is set.
    public string? ServiceAccountJson { get; set; }

    public string? ServiceAccountFile { get; set; }

    // Render-friendly runtime secret names. Values are never sent to frontend.
    public string? AdminProjectId { get; set; }

    public string? AdminClientEmail { get; set; }

    public string? AdminPrivateKey { get; set; }

    public bool IsConfigured =>
        !string.IsNullOrWhiteSpace(ProjectId) &&
        string.Equals(Issuer, ExpectedIssuer, StringComparison.Ordinal);

    public string ExpectedIssuer =>
        string.IsNullOrWhiteSpace(ProjectId)
            ? string.Empty
            : $"https://securetoken.google.com/{ProjectId}";

    public static IReadOnlyList<string> GetValidationErrors(FirebaseAuthenticationOptions options)
    {
        ArgumentNullException.ThrowIfNull(options);

        var errors = new List<string>();
        var hasProject = !string.IsNullOrWhiteSpace(options.ProjectId);
        var hasIssuer = !string.IsNullOrWhiteSpace(options.Issuer);
        var hasCredentialJson = !string.IsNullOrWhiteSpace(options.ServiceAccountJson);
        var hasCredentialFile = !string.IsNullOrWhiteSpace(options.ServiceAccountFile);
        var adminCredentialValues = new[]
        {
            options.AdminProjectId,
            options.AdminClientEmail,
            options.AdminPrivateKey
        };
        var hasAnyAdminCredential = adminCredentialValues.Any(value => !string.IsNullOrWhiteSpace(value));
        var hasCompleteAdminCredential = adminCredentialValues.All(value => !string.IsNullOrWhiteSpace(value));

        if (hasProject != hasIssuer)
        {
            errors.Add("Firebase project ID and issuer must be configured together.");
        }
        else if (hasProject && !string.Equals(options.Issuer, options.ExpectedIssuer,
                     StringComparison.Ordinal))
        {
            errors.Add("Firebase issuer must match the configured project ID.");
        }

        if (hasCredentialJson && hasCredentialFile)
        {
            errors.Add("Firebase service account JSON and file cannot both be configured.");
        }

        if (hasAnyAdminCredential && !hasCompleteAdminCredential)
        {
            errors.Add("Firebase Admin project ID, client email and private key must be configured together.");
        }

        if (hasAnyAdminCredential && (hasCredentialJson || hasCredentialFile))
        {
            errors.Add("Firebase Admin field credentials cannot be combined with service account JSON or file.");
        }

        return errors;
    }
}
