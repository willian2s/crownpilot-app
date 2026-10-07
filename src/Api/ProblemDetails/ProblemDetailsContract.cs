using Microsoft.AspNetCore.Mvc;

namespace CrownPilot.Api.ProblemDetails;

public static class ProblemDetailsCodes
{
    public const string InvalidRequest = "invalid_request";
    public const string AuthenticationRequired = "authentication_required";
    public const string AuthorizationForbidden = "authorization_forbidden";
    public const string ResourceNotFound = "resource_not_found";
    public const string InvalidPlayerTag = "invalid_player_tag";
    public const string PlayerNotFound = "player_not_found";
    public const string ProviderRateLimited = "provider_rate_limited";
    public const string ProviderUnavailable = "provider_unavailable";
    public const string ProviderMisconfigured = "provider_misconfigured";
    public const string VersionConflict = "version_conflict";
    public const string ReauthenticationRequired = "reauthentication_required";
    public const string SemanticValidationFailed = "semantic_validation_failed";
    public const string InternalError = "internal_error";
}

public static class ProblemDetailsContract
{
    public const string Type = "https://www.rfc-editor.org/rfc/rfc9457";

    public static void ApplyDefaults(Microsoft.AspNetCore.Mvc.ProblemDetails problemDetails, int statusCode)
    {
        problemDetails.Status ??= statusCode;
        problemDetails.Type ??= Type;
        problemDetails.Title ??= TitleForStatus(statusCode);
        problemDetails.Detail = null;
        problemDetails.Instance = null;

        if (!problemDetails.Extensions.ContainsKey("code"))
        {
            problemDetails.Extensions["code"] = CodeForStatus(statusCode);
        }
    }

    public static string CodeForStatus(int statusCode) => statusCode switch
    {
        StatusCodes.Status400BadRequest => ProblemDetailsCodes.InvalidRequest,
        StatusCodes.Status401Unauthorized => ProblemDetailsCodes.AuthenticationRequired,
        StatusCodes.Status403Forbidden => ProblemDetailsCodes.AuthorizationForbidden,
        StatusCodes.Status404NotFound => ProblemDetailsCodes.ResourceNotFound,
        StatusCodes.Status409Conflict => ProblemDetailsCodes.VersionConflict,
        StatusCodes.Status422UnprocessableEntity => ProblemDetailsCodes.SemanticValidationFailed,
        StatusCodes.Status429TooManyRequests => ProblemDetailsCodes.ProviderRateLimited,
        StatusCodes.Status503ServiceUnavailable => ProblemDetailsCodes.ProviderUnavailable,
        _ => ProblemDetailsCodes.InternalError
    };

    public static string TitleForStatus(int statusCode) => statusCode switch
    {
        StatusCodes.Status400BadRequest => "The request is invalid.",
        StatusCodes.Status401Unauthorized => "Authentication is required.",
        StatusCodes.Status403Forbidden => "The authenticated subject is not authorized.",
        StatusCodes.Status404NotFound => "The requested resource was not found.",
        StatusCodes.Status409Conflict => "The request conflicts with current state.",
        StatusCodes.Status422UnprocessableEntity => "The request failed semantic validation.",
        StatusCodes.Status429TooManyRequests => "The request was rate limited.",
        StatusCodes.Status503ServiceUnavailable => "A required dependency is unavailable.",
        _ => "The request could not be completed."
    };
}
