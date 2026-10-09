using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Design;

namespace CrownPilot.Infrastructure.Persistence;

public sealed class CrownPilotDbContextFactory : IDesignTimeDbContextFactory<CrownPilotDbContext>
{
    public CrownPilotDbContext CreateDbContext(string[] args)
    {
        var configuredConnectionString = GetConnectionArgument(args) ??
            Environment.GetEnvironmentVariable(DatabaseOptions.ConnectionStringEnvironmentVariable) ??
            Environment.GetEnvironmentVariable("DATABASE_CONNECTION_STRING");
        var environment = Environment.GetEnvironmentVariable("CROWNPILOT__ENVIRONMENT") ??
            Environment.GetEnvironmentVariable("ASPNETCORE_ENVIRONMENT");
        var connectionString = DatabaseConnectionPolicy.ResolveForDesignTime(
            configuredConnectionString,
            environment);

        var options = new DbContextOptionsBuilder<CrownPilotDbContext>()
            .UseNpgsql(connectionString, npgsql => npgsql.MigrationsHistoryTable(
                DatabaseOptions.MigrationHistoryTable,
                DatabaseOptions.MigrationHistorySchema))
            .Options;

        return new CrownPilotDbContext(options);
    }

    private static string? GetConnectionArgument(IReadOnlyList<string> args)
    {
        for (var index = 0; index < args.Count - 1; index++)
        {
            if (args[index].Equals("--connection", StringComparison.OrdinalIgnoreCase))
            {
                return args[index + 1];
            }
        }

        return null;
    }
}
