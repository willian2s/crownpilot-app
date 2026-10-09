using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using CrownPilot.Application.Identity;
using CrownPilot.Infrastructure.Authentication;

namespace CrownPilot.Infrastructure;

public static class DependencyInjection
{
    public static IServiceCollection AddInfrastructure(this IServiceCollection services)
    {
        services.AddOptions<FirebaseAuthenticationOptions>()
            .BindConfiguration(FirebaseAuthenticationOptions.SectionName)
            .PostConfigure<IConfiguration>((options, configuration) =>
            {
                options.AdminProjectId = configuration["FIREBASE_ADMIN_PROJECT_ID"];
                options.AdminClientEmail = configuration["FIREBASE_ADMIN_CLIENT_EMAIL"];
                options.AdminPrivateKey = configuration["FIREBASE_ADMIN_PRIVATE_KEY"];
            })
            .Validate(
                options => FirebaseAuthenticationOptions.GetValidationErrors(options).Count == 0,
                "Firebase authentication configuration is invalid.")
            .ValidateOnStart();
        services.AddSingleton<IFirebaseTokenVerifier, FirebaseTokenVerifier>();
        services.AddSingleton<InfrastructureAssemblyMarker>();
        return services;
    }
}
