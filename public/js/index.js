class BaseComponent {
    constructor() {
        this.judetSearchParam = document.getElementById("judetSearchParam")
        this.statusSearchParam = document.getElementById("statusSearchParam")
        this.formaJuridicaSearchParam = document.getElementById("formaJuridicaSearchParam")

        this.initFilters()
    }

    initFilters() {
        this.query = document.getElementById("searchInput")
        this.county = document.getElementById("countyFilter")
        this.status = document.getElementById("statusFilter")
        this.formaJuridica = document.getElementById("formaFilter")
    }

	getFilters() {
		const params = new URLSearchParams();

		let nume_partial = this.query?.value?.toLowerCase().trim()
		if (nume_partial) {
			params.set('nume_partial', nume_partial)
		}

        if (this.county.value) {
            params.set('judet', this.county.value)
        }

        // always set status if exists since it's defaulted to active
        if (this.status) {
            params.set('status', this.status.value)
        }

        if (this.formaJuridica.value) {
            params.set('forma_juridica', this.formaJuridica.value)
        }

		return params
	}

	setFilters() {
		const params = new URLSearchParams(window.location.search);
		
		if (this.query) {
			this.query.value = params.get("nume_partial") ?? this.query.value
		}

        this.county.value = params.get("judet") ?? this.county.value
        this.status.value = params.get("status") ?? this.status.value
        this.formaJuridica.value = params.get("forma_juridica") ?? this.formaJuridica.value

        return params
	}

    setSearchAction() {
        // TODO: find a better way
        if (this.judetSearchParam) {
            let params = this.getFilters()

            let judet = params.get("judet")
            if (judet) {
                this.judetSearchParam.value = judet
                this.judetSearchParam.disabled = false
            }

            this.statusSearchParam.value = params.get("status")

            let formaJuridica = params.get("forma_juridica")
            if (formaJuridica) {
                this.formaJuridicaSearchParam.value = formaJuridica
                this.formaJuridicaSearchParam.disabled = false
            }
        }
    }

	onChangeFilter() {
        this.setSearchAction()
    }

    Start() {
        this.setFilters()
        this.setSearchAction()
    }
}

function CreateComponent() {
    return new BaseComponent()
}
