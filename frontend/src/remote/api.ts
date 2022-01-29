import axios from "axios";

export default {
  fetchMocks() {
    return axios.get(`/api/mocks/`).then(response => {
      return response.data
    })
  },

  createMocks(data: any) {
    return axios({
      method: 'post',
      url: `/api/mocks/`,
      data: JSON.stringify(data),
      headers: {'Content-Type': 'application/json' }
    }).then(response => {
      return response.data
    })
  },

  deleteMock(id: number) {
    return axios.delete(`/api/mocks/${id}`)
  },
}