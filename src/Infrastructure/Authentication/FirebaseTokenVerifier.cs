using FirebaseAdmin;
using FirebaseAdmin.Auth;
using Google.Apis.Auth.OAuth2;
using CrownPilot.Application.Identity;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Options;

namespace CrownPilot.Infrastructure.Authentication;

/// Firebase Admin SDK adapter. Signature, key discovery/rotation, temporal claims,
/// issuer and audience are verified by the official SDK; this adapter only applies
/// the environment allowlist and exposes the verified UID to the API boundary.
public sealed class FirebaseTokenVerifier : IFirebaseTokenVerifier, IDisposable
{
    private readonly FirebaseAuthenticationOptions options;
    private readonly ILogger<FirebaseTokenVerifier> logger;
    private readonly Lazy<FirebaseAuth?> firebaseAuth;
    private FirebaseApp? firebaseApp;

    public FirebaseTokenVerifier(
        IOptions<FirebaseAuthenticationOptions> options,
        ILogger<FirebaseTokenVerifier> logger)
    {
        this.options = options.Value;
        this.logger = logger;
        firebaseAuth = new Lazy<FirebaseAuth?>(CreateFirebaseAuth);
    }

    public async Task<FirebaseVerifiedToken?> VerifyAsync(
        string idToken,
        CancellationToken cancellationToken = default)
    {
        if (string.IsNullOrWhiteSpace(idToken) || !options.IsConfigured)
        {
            return null;
        }

        try
        {
            var auth = firebaseAuth.Value;
            if (auth is null)
            {
                return null;
            }

            var token = await auth.VerifyIdTokenAsync(idToken, cancellationToken);

            // FirebaseToken is produced only after the SDK has verified the JWT. These
            // comparisons bind that verified result to this runtime environment without
            // reimplementing JWT parsing or cryptography.
            if (string.IsNullOrWhiteSpace(token.Uid) ||
                !string.Equals(token.Uid, token.Subject, StringComparison.Ordinal) ||
                !string.Equals(token.Issuer, options.Issuer, StringComparison.Ordinal) ||
                !string.Equals(token.Audience, options.ProjectId, StringComparison.Ordinal))
            {
                return null;
            }

            return new FirebaseVerifiedToken(token.Uid, token.Issuer, token.Audience);
        }
        catch (OperationCanceledException) when (cancellationToken.IsCancellationRequested)
        {
            throw;
        }
        catch (Exception exception)
        {
            // Verification failures are intentionally indistinguishable at the HTTP
            // boundary. Do not log token, UID, claims, or provider response details.
            logger.LogDebug("Firebase ID token verification failed: {ExceptionType}",
                exception.GetType().Name);
            return null;
        }
    }

    public void Dispose()
    {
        if (firebaseApp is not null)
        {
            firebaseApp.Delete();
            firebaseApp = null;
        }
    }

    private FirebaseAuth? CreateFirebaseAuth()
    {
        if (!options.IsConfigured)
        {
            return null;
        }

        var appOptions = new AppOptions
        {
            Credential = CreateCredential(),
            ProjectId = options.ProjectId
        };

        firebaseApp = FirebaseApp.Create(
            appOptions,
            $"crownpilot-{Guid.NewGuid():N}");

        return FirebaseAuth.GetAuth(firebaseApp);
    }

    private GoogleCredential CreateCredential()
    {
        if (!string.IsNullOrWhiteSpace(options.ServiceAccountJson))
        {
            return CredentialFactory.FromJson<ServiceAccountCredential>(options.ServiceAccountJson)
                .ToGoogleCredential();
        }

        if (!string.IsNullOrWhiteSpace(options.ServiceAccountFile))
        {
            return CredentialFactory.FromFile<ServiceAccountCredential>(options.ServiceAccountFile)
                .ToGoogleCredential();
        }

        if (!string.IsNullOrWhiteSpace(options.AdminProjectId) &&
            !string.IsNullOrWhiteSpace(options.AdminClientEmail) &&
            !string.IsNullOrWhiteSpace(options.AdminPrivateKey))
        {
            var privateKey = options.AdminPrivateKey.Replace("\\n", "\n", StringComparison.Ordinal);
            var initializer = new ServiceAccountCredential.Initializer(options.AdminClientEmail)
            {
                ProjectId = options.AdminProjectId
            }
            .FromPrivateKey(privateKey);
            return new ServiceAccountCredential(initializer).ToGoogleCredential();
        }

        return GoogleCredential.GetApplicationDefault();
    }
}
