namespace CrownPilot.Domain.Identity;

/// Stable internal identity. It is deliberately independent from external provider subjects.
public readonly record struct CrownPilotUserId(Guid Value)
{
    public static CrownPilotUserId New() => new(Guid.NewGuid());
}
