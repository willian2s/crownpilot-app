using CrownPilot.Domain.Identity;

namespace CrownPilot.Application.Identity;

/// External subject produced by authentication after provider claims are verified.
public sealed record AuthenticatedSubject(string FirebaseUid);

/// Application context uses the internal CrownPilot identity, never the external UID directly.
public sealed record AuthContext(CrownPilotUserId CrownPilotUserId);
