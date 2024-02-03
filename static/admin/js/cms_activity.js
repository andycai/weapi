
Alpine.store('activity', {
    init() {
        let currentObj = Alpine.store('current')

        currentObj.prepareResult = (rows, total) => {
            if (!rows) {
                return
            }
            rows.forEach(row => {
                row.view_on_site = currentObj.buildApiUrl(row)
            })
        }
        Alpine.store('queryresult').refresh()
    },
})
