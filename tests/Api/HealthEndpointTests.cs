using System.Net;
using System.Net.Http.Json;
using System.Text.Json;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Mvc.Testing;

namespace CrownPilot.ApiTests;

public sealed class HealthEndpointTests : IClassFixture<WebApplicationFactory<Program>>
{
    private readonly HttpClient client;

    public HealthEndpointTests(WebApplicationFactory<Program> factory)
    {
        client = factory.CreateClient();
    }

    [Fact]
    public async Task LiveHealthDoesNotRequireExternalProvider()
    {
        using var response = await client.GetAsync("/health/live");

        Assert.Equal(HttpStatusCode.OK, response.StatusCode);
        Assert.Equal("Healthy", await response.Content.ReadAsStringAsync());
    }

    [Fact]
    public async Task OpenApiIsGeneratedByApiCode()
    {
        using var response = await client.GetAsync("/openapi/v1.json");
        using var document = await JsonDocument.ParseAsync(await response.Content.ReadAsStreamAsync());

        Assert.Equal(HttpStatusCode.OK, response.StatusCode);
        Assert.StartsWith("3.", document.RootElement.GetProperty("openapi").GetString());
        Assert.True(document.RootElement.GetProperty("paths").TryGetProperty("/api/v1/bootstrap", out _));
    }

    [Fact]
    public async Task SwaggerUiServesGeneratedContractInDevelopment()
    {
        using var response = await client.GetAsync("/docs/index.html");
        var body = await response.Content.ReadAsStringAsync();

        Assert.Equal(HttpStatusCode.OK, response.StatusCode);
        Assert.Contains("id=\"swagger-ui\"", body, StringComparison.Ordinal);
        Assert.Contains("CrownPilot API documentation", body, StringComparison.Ordinal);
    }

    [Fact]
    public async Task SwaggerUiServesGeneratedContractInStaging()
    {
        await using var factory = new WebApplicationFactory<Program>()
            .WithWebHostBuilder(builder =>
                builder.UseSetting(WebHostDefaults.EnvironmentKey, "Staging"));
        using var stagingClient = factory.CreateClient();

        using var response = await stagingClient.GetAsync("/docs/index.html");
        using var openApiResponse = await stagingClient.GetAsync("/openapi/v1.json");

        Assert.Equal(HttpStatusCode.OK, response.StatusCode);
        Assert.Equal(HttpStatusCode.OK, openApiResponse.StatusCode);
    }

    [Fact]
    public async Task ProductionDoesNotExposeDocumentationRoutes()
    {
        await using var factory = new WebApplicationFactory<Program>()
            .WithWebHostBuilder(builder =>
                builder.UseSetting(WebHostDefaults.EnvironmentKey, "Production"));
        using var productionClient = factory.CreateClient();

        using var openApiResponse = await productionClient.GetAsync("/openapi/v1.json");
        using var docsResponse = await productionClient.GetAsync("/docs/index.html");

        Assert.Equal(HttpStatusCode.NotFound, openApiResponse.StatusCode);
        Assert.Equal(HttpStatusCode.NotFound, docsResponse.StatusCode);
    }

    [Fact]
    public async Task MissingRouteUsesProblemDetailsWithoutInternalDetails()
    {
        using var response = await client.GetAsync("/not-found");
        var problem = await response.Content.ReadFromJsonAsync<JsonElement>();

        Assert.Equal(HttpStatusCode.NotFound, response.StatusCode);
        Assert.Equal("The request could not be completed.", problem.GetProperty("title").GetString());
        Assert.True(problem.TryGetProperty("traceId", out _));
        Assert.DoesNotContain("ConnectionStrings", problem.ToString(), StringComparison.OrdinalIgnoreCase);
    }
}
