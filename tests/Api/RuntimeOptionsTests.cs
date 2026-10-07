using CrownPilot.Api.Configuration;

namespace Api.Tests;

public sealed class RuntimeOptionsTests
{
    [Fact]
    public void PreviewRejectsOpenApiExposure()
    {
        var options = new RuntimeOptions
        {
            Environment = "Preview",
            Documentation = new DocumentationOptions { ExposeJson = true, ExposeUi = true },
            Cors = new CorsOptions()
        };

        var errors = RuntimeOptionsValidator.GetValidationErrors(options);

        Assert.Contains(errors, error => error.Contains("OpenAPI exposure", StringComparison.Ordinal));
    }

    [Fact]
    public void NonLocalCorsRequiresHttpsAndRejectsWildcard()
    {
        var options = new RuntimeOptions
        {
            Environment = "Staging",
            Cors = new CorsOptions
            {
                AllowedOrigins = ["*", "http://staging.example"]
            },
            Documentation = new DocumentationOptions { ExposeJson = true, ExposeUi = true },
            Authentication = new AuthenticationBoundaryOptions
            {
                Mode = AuthenticationBoundaryOptions.UnconfiguredMode,
                ContractFixturesEnabled = false
            }
        };

        var errors = RuntimeOptionsValidator.GetValidationErrors(options);

        Assert.Equal(2, errors.Count(error => error.Contains("invalid origin", StringComparison.Ordinal)));
    }

    [Fact]
    public void LocalDefaultsExposeDocsAndUseFixturesOnlyForLocal()
    {
        var options = new RuntimeOptions();

        options.ResolveForHost("Development");

        Assert.Equal(RuntimeEnvironment.Local, options.ResolvedEnvironment);
        Assert.True(options.ExposeOpenApiJson);
        Assert.True(options.ExposeOpenApiUi);
        Assert.True(options.UseContractAuthentication);
    }

    [Fact]
    public void RuntimeEnvironmentMustMatchAspNetEnvironment()
    {
        var options = new RuntimeOptions { Environment = "Preview" };

        var exception = Assert.Throws<InvalidOperationException>(
            () => options.ResolveForHost("Staging"));

        Assert.Contains("must match", exception.Message, StringComparison.Ordinal);
    }

    [Fact]
    public void ContractFixturesAreRejectedOutsideLocal()
    {
        var options = new RuntimeOptions
        {
            Environment = "Staging",
            Authentication = new AuthenticationBoundaryOptions
            {
                Mode = AuthenticationBoundaryOptions.ContractFixtureMode,
                ContractFixturesEnabled = true
            }
        };

        var errors = RuntimeOptionsValidator.GetValidationErrors(options);

        Assert.Contains(errors, error => error.Contains("only in Local", StringComparison.Ordinal));
    }

    [Fact]
    public void FutureProviderModesRemainConfigurationCompatible()
    {
        var options = new RuntimeOptions
        {
            Environment = "Staging",
            Authentication = new AuthenticationBoundaryOptions
            {
                Mode = "Firebase",
                ContractFixturesEnabled = false
            }
        };

        var errors = RuntimeOptionsValidator.GetValidationErrors(options);

        Assert.Empty(errors);
    }
}
