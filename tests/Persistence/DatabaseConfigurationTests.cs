using CrownPilot.Infrastructure.Persistence;

namespace Persistence.Tests;

public sealed class DatabaseConfigurationTests
{
    [Fact]
    [Trait("Category", "Persistence")]
    public void LocalDefaultIsAcceptedWithoutRemoteConfiguration()
    {
        var errors = DatabaseConnectionPolicy.GetValidationErrors(null, "Development");

        Assert.Empty(errors);
    }

    [Fact]
    [Trait("Category", "Persistence")]
    public void PreviewRejectsDatabaseConfiguration()
    {
        var errors = DatabaseConnectionPolicy.GetValidationErrors(
            DatabaseOptions.LocalConnectionString,
            "Preview");

        Assert.Contains(errors, error => error.Contains("Preview", StringComparison.Ordinal));
    }

    [Fact]
    [Trait("Category", "Persistence")]
    public void StagingRequiresCertificateValidationAndNonLoopbackTarget()
    {
        var errors = DatabaseConnectionPolicy.GetValidationErrors(
            "Host=127.0.0.1;Port=5432;Database=crownpilot;Username=runtime;Password=placeholder;SSL Mode=Disable",
            "Staging",
            allowLocalDefault: false);

        Assert.Contains(errors, error => error.Contains("cannot target loopback", StringComparison.Ordinal));
        Assert.Contains(errors, error => error.Contains("VerifyCA or VerifyFull", StringComparison.Ordinal));
    }

    [Fact]
    [Trait("Category", "Persistence")]
    public void StagingRejectsTrustedServerCertificate()
    {
        var errors = DatabaseConnectionPolicy.GetValidationErrors(
            "Host=db.example;Port=5432;Database=crownpilot;Username=runtime;Password=placeholder;SSL Mode=VerifyFull;Trust Server Certificate=true",
            "Staging",
            allowLocalDefault: false);

        Assert.Contains(errors, error => error.Contains("cannot trust", StringComparison.Ordinal));
    }

    [Fact]
    [Trait("Category", "Persistence")]
    public void StagingAcceptsVerifiedTlsReference()
    {
        var errors = DatabaseConnectionPolicy.GetValidationErrors(
            "Host=db.example;Port=5432;Database=crownpilot;Username=runtime;Password=placeholder;SSL Mode=VerifyFull",
            "Staging",
            allowLocalDefault: false);

        Assert.Empty(errors);
    }

    [Fact]
    [Trait("Category", "Persistence")]
    public void ProductionRequiresConnectionReference()
    {
        var errors = DatabaseConnectionPolicy.GetValidationErrors(null, "Production");

        Assert.Contains(errors, error => error.Contains("required", StringComparison.Ordinal));
    }

    [Fact]
    [Trait("Category", "Rls")]
    public void LocalConnectionUsesDedicatedDatabaseAndSchemaConstants()
    {
        Assert.Contains("Database=crownpilot_local", DatabaseOptions.LocalConnectionString,
            StringComparison.Ordinal);
        Assert.Equal("crownpilot", DatabaseOptions.Schema);
        Assert.Equal("public", DatabaseOptions.MigrationHistorySchema);
    }
}
