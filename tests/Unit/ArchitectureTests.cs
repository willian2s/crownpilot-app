using CrownPilot.Domain;

namespace CrownPilot.UnitTests;

public sealed class ArchitectureTests
{
    [Fact]
    public void DomainAssemblyHasNoProviderOrWebReferences()
    {
        var references = typeof(DomainAssemblyMarker).Assembly
            .GetReferencedAssemblies()
            .Select(reference => reference.Name)
            .Where(name => name is not null)
            .ToArray();

        Assert.DoesNotContain(references, name => name!.StartsWith("Microsoft.AspNetCore", StringComparison.Ordinal));
        Assert.DoesNotContain(references, name => name!.StartsWith("Microsoft.EntityFrameworkCore", StringComparison.Ordinal));
        Assert.DoesNotContain(references, name => name!.StartsWith("Npgsql", StringComparison.Ordinal));
        Assert.DoesNotContain(references, name => name!.StartsWith("Firebase", StringComparison.OrdinalIgnoreCase));
        Assert.DoesNotContain(references, name => name!.StartsWith("Supabase", StringComparison.OrdinalIgnoreCase));
    }
}
