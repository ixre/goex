package report

import (
	"bytes"
)

// 生成WEB导出勾选项及脚本
func BuildWebExportCheckOptions(p IDataExportPortal, token string) string {
	portal := p.(*ExportItem)
	buf := bytes.NewBufferString("")
	// 输出Javascript支持库
	buf.WriteString(`<script type="text/javascript">
        var wbExpo={
            //处理请求的Url地址
            urlHandler:'processExport',
            form:null,
            chkInit:function(){
               if(this.form != null)return false;
               this.form = document.forms["export_form"];
               this.form.setAttribute("action",this.urlHandler);
               document.getElementById("params").value = this.getParams();
            },
            getParams:function(){
               var regMatch=/(\?|&)params=(.+)&*/i.exec(location.search);
               return regMatch?regMatch[2]:'';
            },
            submit:function(e){
                this.chkInit();
                this.form.submit();
            }
        };
        </script>`)

	// 输出Wrapper
	buf.WriteString(`<div class="expo-wrapper" id="expo-wrapper">
		<form name="export_form" method="POST" target="export_frame">`)
	// portal
	buf.WriteString("\n<input type=\"hidden\" name=\"portal\" value=\"")
	buf.WriteString(portal.PortalKey)
	buf.WriteString("\"/>\n")
	// params
	buf.WriteString(`<input type="hidden" name="params" value="" id="params"/>`)
	// token
	buf.WriteString("\n<input type=\"hidden\" name=\"token\" value=\"")
	buf.WriteString(token)
	buf.WriteString("\"/>\n")

	// 输出导出格式
	buf.WriteString(`
        <div><strong>选择导出格式</strong></div>
            <ul class="columnList">
            <li class="wbExpo_format_excel"><input type="radio" name="export_format" style="border:none"
               value="excel" checked="checked" id="wbExpo_format_excel"/>
                <label for="wbExpo_format_excel">Excel文件</label>
            </li>
            <li class="wbExpo_format_csv"><input type="radio" name="export_format" style="border:none"
              value="csv" id="wbExpo_format_csv"/>
                <label for="wbExpo_format_csv">CSV数据文件</label>
            </li>
            <li class="wbExpo_format_txt"><input type="radio" name="export_format" style="border:none"
              value="txt" id="wbExpo_format_txt"/>
                <label for="wbExpo_format_txt">文本</label>
            </li>
            </ul>
        <div style="clear:both"></div><br />`)
	// 输出勾选框
	//if e.sqlConfig.ColumnMapping
	colNames := portal.GetColumnNames()
	if len(colNames) == 0 {
		buf.WriteString("<div class=\"expo-no-field\">该导出方案不包含可选择的导出列</div>")
	} else {
		buf.WriteString(`<div class="selColumn"><strong>请选择要导出的列:</strong>
            <ul class="columnList">`)

		for _, col := range colNames {
			buf.WriteString("<li><input type=\"checkbox\" style=\"border:none\" checked=\"checked\"")
			buf.WriteString(" name=\"export_field\" id=\"export_field_")
			buf.WriteString(col.Field)
			buf.WriteString("\" value=\"")
			buf.WriteString(col.Field)
			buf.WriteString("\"/><label for=\"export_field_")
			buf.WriteString(col.Field)
			buf.WriteString("\">")
			buf.WriteString(col.Name)
			buf.WriteString("</label></li>")
		}
		buf.WriteString("</ul></div>")
	}

	buf.WriteString(`<iframe id="export_frame" name="export_frame" style="display:none"></iframe>
        <div style="clear:both"></div><input type="button" class="gra-btn btn-export" onclick="wbExpo.submit()"
         value=" 导出 "/></form></div>`)

	return buf.String()

}
