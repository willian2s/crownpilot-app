namespace CrownPilot.Application.Identity;

/// Result exposed by the Infrastructure Firebase adapter after official SDK verification.
/// Provider claims remain outside the authorization model; only this verified subject crosses
/// the authentication boundary.
public sealed record FirebaseVerifiedToken(
    string Uid,
    string Issuer,
    string Audience);

public interface IFirebaseTokenVerifier
{
    Task<FirebaseVerifiedToken?> VerifyAsync(
        string idToken,
        CancellationToken cancellationToken = default);
}
