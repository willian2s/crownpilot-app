using Microsoft.EntityFrameworkCore;

namespace CrownPilot.Infrastructure.Persistence;

public sealed class CrownPilotDbContext(DbContextOptions<CrownPilotDbContext> options)
    : DbContext(options)
{
    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        // The context reserves the private schema now; 002-06 adds only its
        // approved identity/link entities and their EF-owned constraints.
        modelBuilder.HasDefaultSchema(DatabaseOptions.Schema);
    }
}
