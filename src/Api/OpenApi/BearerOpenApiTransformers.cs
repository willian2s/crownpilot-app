using CrownPilot.Api.Authentication;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.OpenApi;
using Microsoft.OpenApi;

namespace CrownPilot.Api.OpenApi;

public sealed class BearerOpenApiDocumentTransformer : IOpenApiDocumentTransformer
{
    public Task TransformAsync(
        OpenApiDocument document,
        OpenApiDocumentTransformerContext context,
        CancellationToken cancellationToken)
    {
        var components = document.Components ??= new OpenApiComponents();
        components.SecuritySchemes ??= new Dictionary<string, IOpenApiSecurityScheme>();
        components.SecuritySchemes[ContractAuthenticationDefaults.Scheme] =
            new OpenApiSecurityScheme
            {
                Type = SecuritySchemeType.Http,
                Scheme = "bearer",
                BearerFormat = "Firebase ID Token",
                Description = "Send only a server-verified Firebase ID Token as a bearer token."
            };

        if (components.Schemas is { } schemas &&
            schemas.TryGetValue("ProblemDetails", out var problemDetails) &&
            problemDetails is OpenApiSchema problemDetailsSchema)
        {
            if (problemDetailsSchema.Properties is not { } properties)
            {
                return Task.CompletedTask;
            }

            properties.TryAdd("code", new OpenApiSchema
            {
                Type = JsonSchemaType.String,
                Description = "Stable non-sensitive application error code."
            });
            properties.TryAdd("traceId", new OpenApiSchema
            {
                Type = JsonSchemaType.String,
                Description = "Correlation identifier for support; never a secret."
            });
        }

        return Task.CompletedTask;
    }
}

public sealed class BearerOpenApiOperationTransformer : IOpenApiOperationTransformer
{
    public Task TransformAsync(
        OpenApiOperation operation,
        OpenApiOperationTransformerContext context,
        CancellationToken cancellationToken)
    {
        var requiresAuthorization = context.Description.ActionDescriptor.EndpointMetadata
            .OfType<IAuthorizeData>()
            .Any();

        if (!requiresAuthorization)
        {
            return Task.CompletedTask;
        }

        operation.Security ??= [];
        operation.Security.Add(new OpenApiSecurityRequirement
        {
            [new OpenApiSecuritySchemeReference(
                ContractAuthenticationDefaults.Scheme,
                context.Document,
                externalResource: null)] = []
        });

        var responses = operation.Responses ??= new OpenApiResponses();
        responses.TryAdd("401", new OpenApiResponse
        {
            Description = "Bearer token is absent or invalid."
        });
        responses.TryAdd("403", new OpenApiResponse
        {
            Description = "Authenticated subject is not authorized."
        });

        return Task.CompletedTask;
    }
}
