using System.Diagnostics;
using CrownPilot.Infrastructure;
using Microsoft.AspNetCore.Diagnostics;
using Microsoft.AspNetCore.Diagnostics.HealthChecks;
using Microsoft.AspNetCore.Http.HttpResults;
using Microsoft.AspNetCore.Mvc;

var builder = WebApplication.CreateBuilder(args);

var port = Environment.GetEnvironmentVariable("PORT");
if (int.TryParse(port, out var portNumber) && portNumber is > 0 and <= 65535)
{
    // Hosting platforms commonly provide PORT at runtime; keep local/container
    // development on the ASP.NET Core default when the variable is absent.
    builder.WebHost.UseUrls($"http://0.0.0.0:{portNumber}");
}

builder.Services.AddProblemDetails(options =>
{
    options.CustomizeProblemDetails = context =>
    {
        context.ProblemDetails.Extensions["traceId"] =
            Activity.Current?.Id ?? context.HttpContext.TraceIdentifier;
    };
});
builder.Services.AddHealthChecks();
builder.Services.AddOpenApi("v1");
builder.Services.AddInfrastructure();

var app = builder.Build();

app.UseExceptionHandler();
app.UseStatusCodePages(async statusCodeContext =>
{
    var httpContext = statusCodeContext.HttpContext;
    var response = httpContext.Response;

    if (response.HasStarted || response.ContentLength is not null || response.StatusCode < 400)
    {
        return;
    }

    var problemDetailsService = httpContext.RequestServices.GetRequiredService<IProblemDetailsService>();
    await problemDetailsService.WriteAsync(new ProblemDetailsContext
    {
        HttpContext = httpContext,
        ProblemDetails = new ProblemDetails
        {
            Status = response.StatusCode,
            Title = "The request could not be completed.",
            Type = "https://www.rfc-editor.org/rfc/rfc9457"
        }
    });
});

app.MapHealthChecks("/health/live", new HealthCheckOptions
{
    Predicate = _ => false
});

if (app.Environment.IsDevelopment() || app.Environment.IsEnvironment("Staging"))
{
    app.MapOpenApi("/openapi/{documentName}.json");

    app.UseSwaggerUI(options =>
    {
        options.RoutePrefix = "docs";
        options.SwaggerEndpoint("/openapi/v1.json", "CrownPilot API v1");
        options.DocumentTitle = "CrownPilot API documentation";
    });
}

app.MapGet("/api/v1/bootstrap", () =>
    TypedResults.Ok(new BootstrapStatus("ready", "bootstrap")))
    .WithName("GetBootstrapStatus")
    .WithTags("Bootstrap")
    .WithSummary("Reports that the local application foundation is running.")
    .Produces<BootstrapStatus>(StatusCodes.Status200OK)
    .ProducesProblem(StatusCodes.Status500InternalServerError);

app.Run();

public sealed record BootstrapStatus(string Status, string Stage);

public partial class Program
{
}
