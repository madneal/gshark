import { getDict } from "@/utils/dictionary";
export default {
    data() {
        return {
            page: 1,
            total: 10,
            pageSize: 100,
            tableData: [],
            searchInfo: {}
        }
    },
    methods: {
        handleSizeChange(val) {
            this.pageSize = val
            this.getTableData()
        },
        handleCurrentChange(val) {
            this.page = val;
            this.getTableData()
        },
        async getTableData(page = this.page, pageSize = this.pageSize) {
            const table = await this.listApi({ page, pageSize, ...this.searchInfo })
            if (table.code === 0) {
                this.tableData = table.data.list
                this.total = table.data.total
                this.page = table.data.page
                this.pageSize = table.data.pageSize
            }
        }
    }
}