namespace CrownPilot.Api.Configuration;

public enum RuntimeEnvironment
{
    Local,
    Preview,
    Staging,
    Production
}

public static class RuntimeEnvironmentNames
{
    public static bool TryParse(string? value, out RuntimeEnvironment environment)
    {
        environment = default;

        if (string.IsNullOrWhiteSpace(value))
        {
            return false;
        }

        if (value.Equals("Development", StringComparison.OrdinalIgnoreCase) ||
            value.Equals(nameof(RuntimeEnvironment.Local), StringComparison.OrdinalIgnoreCase))
        {
            environment = RuntimeEnvironment.Local;
            return true;
        }

        return Enum.TryParse(value, ignoreCase: true, out environment) &&
            Enum.IsDefined(environment);
    }

    public static bool TryParseHost(string? hostEnvironment, out RuntimeEnvironment environment) =>
        TryParse(hostEnvironment, out environment);
}
