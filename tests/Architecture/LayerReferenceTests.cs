using System.IO;

namespace CrownPilot.ArchitectureTests;

public sealed class LayerReferenceTests
{
    [Fact]
    public void InfrastructureReferencesApplicationAndDomain()
    {
        var project = ReadProject("src/Infrastructure/Infrastructure.csproj");

        Assert.Contains("..\\Application\\Application.csproj", project, StringComparison.Ordinal);
        Assert.Contains("..\\Domain\\Domain.csproj", project, StringComparison.Ordinal);
    }

    [Fact]
    public void ApiCompositionRootReferencesApplicationAndInfrastructure()
    {
        var project = ReadProject("src/Api/Api.csproj");

        Assert.Contains("..\\Application\\Application.csproj", project, StringComparison.Ordinal);
        Assert.Contains("..\\Infrastructure\\Infrastructure.csproj", project, StringComparison.Ordinal);
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
