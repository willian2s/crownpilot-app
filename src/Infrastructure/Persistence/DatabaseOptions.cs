using Npgsql;

namespace CrownPilot.Infrastructure.Persistence;

public sealed class DatabaseOptions
{
    public const string SectionName = "CrownPilot:Database";
    public const string ConnectionStringEnvironmentVariable =
        "CROWNPILOT__DATABASE__CONNECTIONSTRING";
    public const string Schema = "crownpilot";
    public const string MigrationHistorySchema = "public";
    public const string MigrationHistoryTable = "__EFMigrationsHistory";

    // This credential exists only inside the disposable local PostgreSQL
    // container. Staging and Production receive a different value at runtime.
    public const string LocalConnectionString =
        "Host=127.0.0.1;Port=54322;Database=crownpilot_local;Username=postgres;Password=postgres;SSL Mode=Disable;Application Name=CrownPilot.DbScripts";

    public string? ConnectionString { get; set; }
}

public static class DatabaseConnectionPolicy
{
    public static IReadOnlyList<string> GetValidationErrors(
        string? connectionString,
        string? environment,
        bool allowLocalDefault = true)
    {
        var errors = new List<string>();
        var normalizedEnvironment = NormalizeEnvironment(environment);

        if (normalizedEnvironment is null)
        {
            errors.Add("Database environment must be Local, Preview, Staging, or Production.");
            return errors;
        }

        if (normalizedEnvironment == "Preview")
        {
            if (!string.IsNullOrWhiteSpace(connectionString))
            {
                errors.Add("Preview cannot be configured with a database connection string.");
            }

            return errors;
        }

        if (string.IsNullOrWhiteSpace(connectionString))
        {
            if (normalizedEnvironment == "Local" && allowLocalDefault)
            {
                return errors;
            }

            errors.Add($"A database connection string is required for {normalizedEnvironment}.");
            return errors;
        }

        NpgsqlConnectionStringBuilder parsed;
        try
        {
            parsed = new NpgsqlConnectionStringBuilder(connectionString);
        }
        catch (ArgumentException)
        {
            errors.Add("Database connection string is not a valid Npgsql connection string.");
            return errors;
        }

        if (string.IsNullOrWhiteSpace(parsed.Host) || string.IsNullOrWhiteSpace(parsed.Database) ||
            string.IsNullOrWhiteSpace(parsed.Username))
        {
            errors.Add("Database connection string must define host, database, and username.");
        }

        if (normalizedEnvironment == "Local")
        {
            if (!IsLoopback(parsed.Host ?? string.Empty))
            {
                errors.Add("Local database connection must target loopback PostgreSQL.");
            }
        }
        else
        {
            if (IsLoopback(parsed.Host ?? string.Empty))
            {
                errors.Add($"{normalizedEnvironment} database connection cannot target loopback PostgreSQL.");
            }

            if (parsed.SslMode is not (SslMode.VerifyCA or SslMode.VerifyFull))
            {
                errors.Add($"{normalizedEnvironment} database connection must use VerifyCA or VerifyFull.");
            }

            if (RequestsUnverifiedCertificate(connectionString))
            {
                errors.Add($"{normalizedEnvironment} database connection cannot trust an unverified server certificate.");
            }
        }

        return errors;
    }

    public static string ResolveForDesignTime(string? configuredConnectionString, string? environment)
    {
        var normalizedEnvironment = NormalizeEnvironment(environment);
        var connectionString = string.IsNullOrWhiteSpace(configuredConnectionString)
            ? normalizedEnvironment == "Local"
                ? DatabaseOptions.LocalConnectionString
                : null
            : configuredConnectionString;

        var errors = GetValidationErrors(connectionString, normalizedEnvironment, allowLocalDefault: false);
        if (errors.Count > 0)
        {
            throw new InvalidOperationException(string.Join(" ", errors));
        }

        return connectionString!;
    }

    public static string? NormalizeEnvironment(string? value)
    {
        if (string.IsNullOrWhiteSpace(value) ||
            value.Equals("Development", StringComparison.OrdinalIgnoreCase) ||
            value.Equals("Local", StringComparison.OrdinalIgnoreCase))
        {
            return "Local";
        }

        if (value.Equals("Preview", StringComparison.OrdinalIgnoreCase))
        {
            return "Preview";
        }

        if (value.Equals("Staging", StringComparison.OrdinalIgnoreCase))
        {
            return "Staging";
        }

        if (value.Equals("Production", StringComparison.OrdinalIgnoreCase))
        {
            return "Production";
        }

        return null;
    }

    private static bool IsLoopback(string host) =>
        host.Equals("localhost", StringComparison.OrdinalIgnoreCase) ||
        host.Equals("127.0.0.1", StringComparison.OrdinalIgnoreCase) ||
        host.Equals("::1", StringComparison.OrdinalIgnoreCase) ||
        host.Equals("[::1]", StringComparison.OrdinalIgnoreCase);

    private static bool RequestsUnverifiedCertificate(string connectionString)
    {
        foreach (var part in connectionString.Split(';', StringSplitOptions.RemoveEmptyEntries))
        {
            var separator = part.IndexOf('=');
            if (separator < 1)
            {
                continue;
            }

            var key = part[..separator].Trim()
                .Replace(" ", string.Empty, StringComparison.Ordinal)
                .Replace("-", string.Empty, StringComparison.Ordinal)
                .Replace("_", string.Empty, StringComparison.Ordinal);
            var value = part[(separator + 1)..].Trim().Trim('"', '\'');

            if (key.Equals("trustservercertificate", StringComparison.OrdinalIgnoreCase) &&
                bool.TryParse(value, out var trustsCertificate) && trustsCertificate)
            {
                return true;
            }
        }

        return false;
    }
}
