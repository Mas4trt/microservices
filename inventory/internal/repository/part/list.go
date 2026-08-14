package part

import (
	"context"

	"github.com/Mas4trt/microservices/inventory/internal/model"
	"github.com/Mas4trt/microservices/inventory/internal/repository/converter"
	repoConverter "github.com/Mas4trt/microservices/inventory/internal/repository/converter"
	repoModel "github.com/Mas4trt/microservices/inventory/internal/repository/model"
)

func (r *repository) List(_ context.Context, filter model.PartsFilter) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	convFilter := repoConverter.FilterToRepoModel(filter)

	uuids := stringValueSet(convFilter.UUIDs)
	names := stringValueSet(convFilter.Names)
	countries := stringValueSet(convFilter.ManufacturerCountries)
	tags := stringValueSet(convFilter.Tags)
	categories := categorySet(convFilter.Categories)

	result := make([]repoModel.Part, 0, len(r.parts))

	for _, part := range r.parts {
		if !matchesSet(uuids, part.Uuid) {
			continue
		}
		if !matchesSet(names, part.Name) {
			continue
		}
		if !matchesCategorySet(categories, part.Category) {
			continue
		}
		if !matchesSet(countries, part.Manufacturer.Country) {
			continue
		}
		if !matchesAnySet(tags, part.Tags) {
			continue
		}

		result = append(result, part)
	}

	return converter.PartsToModel(result), nil
}

func stringValueSet(values []string) map[string]struct{} {
	if len(values) == 0 {
		return nil
	}

	set := make(map[string]struct{}, len(values))

	for _, v := range values {
		set[v] = struct{}{}
	}

	return set
}

func categorySet(values []repoModel.Category) map[repoModel.Category]struct{} {
	if len(values) == 0 {
		return nil
	}

	set := make(map[repoModel.Category]struct{}, len(values))
	for _, v := range values {
		set[v] = struct{}{}
	}

	return set
}

// matchesSet: nil/пустое множество означает "фильтр не задан" -> всегда true.
func matchesSet(set map[string]struct{}, value string) bool {
	if set == nil {
		return true
	}

	_, ok := set[value]
	return ok
}

func matchesCategorySet(set map[repoModel.Category]struct{}, value repoModel.Category) bool {
	if set == nil {
		return true
	}

	_, ok := set[value]
	return ok
}

// matchesAnySet: ИЛИ внутри поля — хотя бы один тег детали входит в фильтр.
func matchesAnySet(set map[string]struct{}, values []string) bool {
	if set == nil {
		return true
	}

	for _, v := range values {
		if _, ok := set[v]; ok {
			return true
		}
	}

	return false
}
