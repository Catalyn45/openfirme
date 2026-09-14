class HomePage extends SearchPage {
    init() {
        this.initSearchBar()
        this.initFilters()
    }

	onChangeFilter() { }

    Start() { }
}

function CreateComponent() {
    return new HomePage()
}
