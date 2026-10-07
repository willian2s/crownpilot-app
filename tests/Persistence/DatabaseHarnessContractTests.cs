using System.Text.Json;

namespace Persistence.Tests;

public sealed class DatabaseHarnessContractTests
{
    [Fact]
    [Trait("Category", "Persistence")]
    public void SecuritySqlDoesNotOwnDomainTablesOrPolicies()
    {
        var sql = ReadRepositoryFile("database/security/001-roles-schema.sql");
        var normalized = sql.ToLowerInvariant();

        Assert.DoesNotContain("create table", normalized, StringComparison.Ordinal);
        Assert.DoesNotContain("create policy", normalized, StringComparison.Ordinal);
        Assert.DoesNotContain("create schema", normalized, StringComparison.Ordinal);
        Assert.Contains("create role", normalized, StringComparison.Ordinal);
        Assert.Contains("alter default privileges", normalized, StringComparison.Ordinal);
    }

    [Fact]
    [Trait("Category", "Rls")]
    public void ComposeUsesPinnedDisposablePostgresAndNoPersistentVolume()
    {
        var compose = ReadRepositoryFile("database/docker-compose.yml");

        Assert.Contains("postgres:17.6-alpine", compose, StringComparison.Ordinal);
        Assert.Contains("/var/lib/postgresql/data", compose, StringComparison.Ordinal);
        Assert.DoesNotContain("volumes:", compose, StringComparison.Ordinal);
    }

    [Fact]
    [Trait("Category", "Persistence")]
    public void DatabaseScriptsRemainPublicNpmInterface()
    {
        using var document = JsonDocument.Parse(ReadRepositoryFile("frontend/package.json"));
        var scripts = document.RootElement.GetProperty("scripts");

        foreach (var operation in new[] { "start", "stop", "migrate", "security", "fixtures", "reset" })
        {
            var script = scripts.GetProperty($"db:{operation}").GetString();
            Assert.Contains("scripts/db-gate.mjs", script, StringComparison.Ordinal);
        }
    }

    private static string ReadRepositoryFile(string relativePath)
    {
        var directory = new DirectoryInfo(AppContext.BaseDirectory);
        while (directory is not null && !File.Exists(Path.Combine(directory.FullName, "CrownPilot.sln")))
        {
            directory = directory.Parent;
        }

        Assert.NotNull(directory);
        return File.ReadAllText(Path.Combine(directory!.FullName, relativePath));
    }
}
