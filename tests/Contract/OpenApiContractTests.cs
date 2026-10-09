using System.Net;
using System.Text.Json;
using Microsoft.AspNetCore.Mvc.Testing;

namespace Contract.Tests;

public sealed class OpenApiContractTests(WebApplicationFactory<Program> factory)
    : IClassFixture<WebApplicationFactory<Program>>
{
    private readonly HttpClient client = factory.CreateClient();

    [Fact]
    public async Task BootstrapContractContainsOnlyDocumentedFoundationRoute()
    {
        using var response = await client.GetAsync("/openapi/v1.json");
        using var document = await JsonDocument.ParseAsync(await response.Content.ReadAsStreamAsync());
        var paths = document.RootElement.GetProperty("paths");

        Assert.Equal(HttpStatusCode.OK, response.StatusCode);
        Assert.True(paths.TryGetProperty("/api/v1/bootstrap", out var bootstrap));
        Assert.True(bootstrap.TryGetProperty("get", out var get));
        Assert.Equal("GetBootstrapStatus", get.GetProperty("operationId").GetString());
        Assert.True(get.GetProperty("security").EnumerateArray().Any());
        Assert.True(get.GetProperty("responses").TryGetProperty("401", out _));
        Assert.True(get.GetProperty("responses").TryGetProperty("403", out _));
        Assert.True(document.RootElement.GetProperty("components")
            .GetProperty("securitySchemes").TryGetProperty("Bearer", out var bearer));
        Assert.Equal("http", bearer.GetProperty("type").GetString());
        Assert.Equal("bearer", bearer.GetProperty("scheme").GetString());
        Assert.False(paths.TryGetProperty("/api/v1/me/player-link", out _));
    }
}
