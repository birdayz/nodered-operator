module.exports = {
  adminAuth: {
    type: "credentials",
    users: [
      {{range .AdminAuth.Users}} 
      {
        "username" : "{{ .Username }}",
        "password" : "{{ .Password }}",
        "permissions" : "{{ .Permissions }}"
      }
      {{end}} 
    ]
  },
  editorTheme: {
    page: {
      title: "{{ .Settings.EditorTheme.Page.Title }}"
    }
  }
}
