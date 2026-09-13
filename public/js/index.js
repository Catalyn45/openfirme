class HomePage extends SearchPage {
    init() {
        this.initSearchBar()
        this.initFilters()
    }

	onChangeFilter() {
    }

    Start() {
        this.setFilters()
    }
}

function CreateComponent() {
    return new HomePage()
}
