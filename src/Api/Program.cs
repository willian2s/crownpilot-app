using System.Diagnostics;
using CrownPilot.Api.Authentication;
using CrownPilot.Api.Configuration;
using CrownPilot.Api.OpenApi;
using CrownPilot.Api.ProblemDetails;
using CrownPilot.Api;
using CrownPilot.Infrastructure;
using CrownPilot.Infrastructure.Persistence;
using Microsoft.AspNetCore.Authentication;
using Microsoft.AspNetCore.Diagnostics;
using Microsoft.AspNetCore.Diagnostics.HealthChecks;
using Microsoft.AspNetCore.Http.HttpResults;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Options;

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
        var statusCode = context.ProblemDetails.Status ?? context.HttpContext.Response.StatusCode;
        ProblemDetailsContract.ApplyDefaults(context.ProblemDetails, statusCode);
        context.ProblemDetails.Extensions["traceId"] =
            Activity.Current?.Id ?? context.HttpContext.TraceIdentifier;
    };
});
builder.Services.AddHealthChecks();
builder.Services.AddCors();
builder.Services.AddOptions<RuntimeOptions>()
    .BindConfiguration(RuntimeOptions.SectionName)
    .Configure<IHostEnvironment>((options, hostEnvironment) =>
        options.ResolveForHost(hostEnvironment.EnvironmentName))
    .ValidateOnStart();
builder.Services.AddSingleton<IValidateOptions<RuntimeOptions>, RuntimeOptionsValidator>();
// Database credentials are bound for migration/persistence composition only;
// this task deliberately does not open a connection or migrate during startup.
builder.Services.AddOptions<DatabaseOptions>()
    .BindConfiguration(DatabaseOptions.SectionName);
builder.Services.AddAuthentication(ContractAuthenticationDefaults.Scheme)
    .AddScheme<AuthenticationSchemeOptions, ContractBearerAuthenticationHandler>(
        ContractAuthenticationDefaults.Scheme,
        _ => { });
builder.Services.AddAuthorization(options =>
{
    options.AddPolicy(ContractAuthenticationDefaults.BootstrapPolicy, policy =>
    {
        policy.AddAuthenticationSchemes(ContractAuthenticationDefaults.Scheme);
        policy.RequireAuthenticatedUser();
        policy.RequireAssertion(context =>
        {
            if (context.Resource is not HttpContext httpContext)
            {
                return false;
            }

            return context.User.HasClaim(
                ContractAuthenticationDefaults.PermissionClaim,
                ContractAuthenticationDefaults.BootstrapReadPermission);
        });
    });
});
builder.Services.AddOpenApi("v1", options =>
{
    options.AddDocumentTransformer<BearerOpenApiDocumentTransformer>();
    options.AddOperationTransformer<BearerOpenApiOperationTransformer>();
});
builder.Services.AddInfrastructure();

var app = builder.Build();
var runtimeOptions = app.Services.GetRequiredService<IOptions<RuntimeOptions>>().Value;

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
            Type = ProblemDetailsContract.Type
        }
    });
});
app.UseCors(runtimeOptions.ConfigureCorsPolicy);
app.UseAuthentication();
app.UseAuthorization();

app.MapHealthChecks("/health/live", new HealthCheckOptions
{
    Predicate = _ => false
});

if (runtimeOptions.ExposeOpenApiJson)
{
    app.MapOpenApi("/openapi/{documentName}.json");
}

if (runtimeOptions.ExposeOpenApiUi)
{
    app.UseSwaggerUI(options =>
    {
        options.RoutePrefix = "docs";
        options.SwaggerEndpoint("/openapi/v1.json", "CrownPilot API v1");
        options.DocumentTitle = "CrownPilot API documentation";
    });
}

app.MapGet("/api/v1/bootstrap", () =>
    TypedResults.Ok(new BootstrapStatus("ready", "bootstrap")))
    .RequireAuthorization(ContractAuthenticationDefaults.BootstrapPolicy)
    .WithName("GetBootstrapStatus")
    .WithTags("Bootstrap")
    .WithSummary("Reports that the authenticated application foundation is running.")
    .Produces<BootstrapStatus>(StatusCodes.Status200OK)
    .ProducesProblem(StatusCodes.Status401Unauthorized)
    .ProducesProblem(StatusCodes.Status403Forbidden)
    .ProducesProblem(StatusCodes.Status500InternalServerError);

app.Run();
