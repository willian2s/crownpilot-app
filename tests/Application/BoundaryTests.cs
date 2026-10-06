using System.IO;

namespace Application.Tests;

public sealed class BoundaryTests
{
    [Fact]
    public void ApplicationDependsOnDomainButNotProviderSdk()
    {
        var project = ReadProject("src/Application/Application.csproj");

        Assert.Contains("..\\Domain\\Domain.csproj", project, StringComparison.Ordinal);
        Assert.DoesNotContain("Firebase", project, StringComparison.OrdinalIgnoreCase);
        Assert.DoesNotContain("Supabase", project, StringComparison.OrdinalIgnoreCase);
        Assert.DoesNotContain("Microsoft.EntityFrameworkCore", project, StringComparison.Ordinal);
        Assert.DoesNotContain("Npgsql", project, StringComparison.Ordinal);
    }

    private static string ReadProject(string relativePath)
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
