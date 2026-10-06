using Microsoft.Extensions.DependencyInjection;

namespace CrownPilot.Infrastructure;

public static class DependencyInjection
{
    public static IServiceCollection AddInfrastructure(this IServiceCollection services)
    {
        // Provider adapters and database options are deliberately added by later subtasks.
        // Registering the marker keeps the composition root explicit without opening a connection.
        services.AddSingleton<InfrastructureAssemblyMarker>();
        return services;
    }
}
