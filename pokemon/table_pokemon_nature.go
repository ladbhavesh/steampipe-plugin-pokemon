package pokemon

import (
	"context"

	"github.com/mtslzr/pokeapi-go"
	"github.com/mtslzr/pokeapi-go/structs"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tablePokemonNature(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "pokemon_nature",
		Description: "Natures influence how a Pokémon's stats grow.",
		List: &plugin.ListConfig{
			Hydrate: listNatures,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"name"}),
			// TODO: Add support for 'id' key column
			//KeyColumns: plugin.AnyColumn([]string{"id", "name"}),
			Hydrate: getNature,
			// Bad error message is a result of https://github.com/mtslzr/pokeapi-go/issues/29
			ShouldIgnoreError: isNotFoundError([]string{"invalid character 'N' looking for beginning of value"}),
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name for this resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "decreased_stat",
				Description: "The stat decreased by 10% in Pokémon with this nature.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNature,
			},
			{
				Name:        "increased_stat",
				Description: "The stat increased by 10% in Pokémon with this nature.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNature,
			},
			{
				Name:        "hates_flavor",
				Description: "The flavor hated by Pokémon with this nature.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNature,
			},
			{
				Name:        "pokeathlon_stat_changes",
				Description: "A list of Pokéathlon stats this nature effects and how much it effects them.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNature,
			},
			{
				Name:        "move_battle_style_preferences",
				Description: "A list of battle styles and how likely a Pokémon with this nature is to use them in the Battle Palace or Battle Tent.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNature,
			},
			{
				Name:        "names",
				Description: "The name of this resource listed in different languages.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getNature,
			},
			{
				Name:        "id",
				Description: "The identifier for this resource.",
				Type:        proto.ColumnType_INT,
				Hydrate:     getNature,
				Transform:   transform.FromGo(),
			},
		},
	}
}

func listNatures(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listNatures")

	offset := 0

	for true {
		resources, err := pokeapi.Resource("nature", offset)

		if err != nil {
			plugin.Logger(ctx).Error("pokemon_nature.listNatures", "query_error", err)
			return nil, err
		}

		for _, pokemon := range resources.Results {
			d.StreamListItem(ctx, pokemon)
		}

		// No next URL returned
		if len(resources.Next) == 0 {
			break
		}

		urlOffset, err := extractUrlOffset(resources.Next)
		if err != nil {
			plugin.Logger(ctx).Error("pokemon_nature.listNatures", "extract_url_offset_error", err)
			return nil, err
		}

		// Set next offset
		offset = urlOffset
	}

	return nil, nil
}

func getNature(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getNature")

	var name string

	if h.Item != nil {
		result := h.Item.(structs.Result)
		name = result.Name
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	logger.Debug("Name", name)

	pokemon, err := pokeapi.Nature(name)

	if err != nil {
		plugin.Logger(ctx).Error("pokemon_nature.pokemonGet", "query_error", err)
		return nil, err
	}

	return pokemon, nil
}
