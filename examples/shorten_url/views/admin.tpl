<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
        "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
    <link rel="stylesheet" href="/static/index.css">
</head>
<body>
<table border="1">
    <tr>
    <th>名字</th>
    <th>地址</th>
    <th>短链</th>
    <th>访问</th>
    </tr>
{{range .Slice}}
    <tr>
        <td>{{.Name}}</td>
        <td class="special"><a href="http://{{localHost}}/t/{{.UniqueId}}"> /t/{{.UniqueId}} </a></td>
        <td>{{.Url}}</td>
        <td>{{.Count}}</td>
    </tr>
{{end}}
</table>

<a href="/dwz">添加链接</a>

</body>
</html>
