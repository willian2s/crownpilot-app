using Microsoft.EntityFrameworkCore.Infrastructure;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace CrownPilot.Infrastructure.Persistence.Migrations;

/// <inheritdoc />
[Migration("20261007000000_PrepareCrownPilotSchema")]
public partial class PrepareCrownPilotSchema : Migration
{
    /// <inheritdoc />
    protected override void Up(MigrationBuilder migrationBuilder)
    {
        migrationBuilder.EnsureSchema(
            name: "crownpilot");
    }

    /// <inheritdoc />
    protected override void Down(MigrationBuilder migrationBuilder)
    {
        migrationBuilder.DropSchema(
            name: "crownpilot");
    }
}
