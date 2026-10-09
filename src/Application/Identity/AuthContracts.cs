using CrownPilot.Domain.Identity;
using System.Collections.Concurrent;

namespace CrownPilot.Application.Identity;

/// External subject produced by authentication after provider claims are verified.
public sealed record AuthenticatedSubject
{
    public AuthenticatedSubject(string firebaseUid)
    {
        if (string.IsNullOrWhiteSpace(firebaseUid))
        {
            throw new ArgumentException("A verified Firebase UID is required.", nameof(firebaseUid));
        }

        FirebaseUid = firebaseUid;
    }

    public string FirebaseUid { get; }
}

/// Application context uses the internal CrownPilot identity, never the external UID directly.
public sealed record AuthContext(CrownPilotUserId CrownPilotUserId);

/// Maps a verified external subject to the provider-independent CrownPilot identity.
///
/// Authentication handlers must not invoke this port. Use cases invoke it when they
/// need a local identity, after authentication has already succeeded.
public interface IUserIdentityResolver
{
    Task<AuthContext> EnsureCrownPilotUserAsync(
        AuthenticatedSubject subject,
        CancellationToken cancellationToken = default);
}

/// Application use-case boundary for the later identity persistence integration.
public sealed class EnsureCrownPilotUser(IUserIdentityResolver resolver)
{
    public Task<AuthContext> ExecuteAsync(
        AuthenticatedSubject subject,
        CancellationToken cancellationToken = default) =>
        resolver.EnsureCrownPilotUserAsync(subject, cancellationToken);
}

/// In-memory contract fake for Application tests. It does not represent persistence.
public sealed class FakeUserIdentityResolver : IUserIdentityResolver
{
    private readonly ConcurrentDictionary<string, CrownPilotUserId> identities = new();

    public int CallCount => Volatile.Read(ref callCount);

    private int callCount;

    public Task<AuthContext> EnsureCrownPilotUserAsync(
        AuthenticatedSubject subject,
        CancellationToken cancellationToken = default)
    {
        ArgumentNullException.ThrowIfNull(subject);
        cancellationToken.ThrowIfCancellationRequested();
        Interlocked.Increment(ref callCount);

        var crownPilotUserId = identities.GetOrAdd(
            subject.FirebaseUid,
            _ => CrownPilotUserId.New());

        return Task.FromResult(new AuthContext(crownPilotUserId));
    }
}
