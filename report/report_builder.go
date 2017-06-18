package report

import (
	"bytes"
	"strconv"
)

// 生成WEB导出勾选项及脚本
func BuildWebExportCheckOptions(p IDataExportPortal) string {
	portal := p.(*ExportItem)
	buf := bytes.NewBufferString("")
	// 输出Javascript支持库
	buf.WriteString(`<script type="text/javascript">
        var wbExpo={
            config:{
                //处理请求的Url地址
                urlHandler:'',
                //处理生成Json的对象
                jsonHandler:null,
                params:null,
                page:null
            },
            doExport:function(e){
                var data=this.config.jsonHandler.toQueryString('expo-wrapper');
                if(this.config.params==null){
                   var regMatch=/(\?|&)params=(.+)&*/i.exec(location.search);
                   this.config.params=regMatch?regMatch[2]:'';
                }
                if(!this.config.page){
                    this.config.page=document.getElementById('export_frame');
                }
                this.config.page.src=this.config.urlHandler
                             + (this.config.urlHandler.indexOf('?')==-1?'?':'&')
                             + 'portal=' + e + '&' + data
                             + '&params=' + this.config.params;
            }
        };
        </script>`)

	// 输出Wrapper
	buf.WriteString(`<div class="expo-wrapper" id="expo-wrapper">`)
	// 输出导出格式
	buf.WriteString(`
        <div><strong>选择导出格式</strong></div>
            <ul class="columnList">
            <li class="wbExpo_format_excel"><input checked="checked" field="export_format" style="border:none" name="wbExpo_format" type="radio" value="excel" id="wbExpo_format_excel"/>
                <label for="wbExpo_format_excel">Excel文件</label>
            </li>
            <li class="wbExpo_format_csv"><input type="radio" field="export_format" style="border:none" name="wbExpo_format" value="csv" id="wbExpo_format_csv"/>
                <label for="wbExpo_format_csv">CSV数据文件</label>
            </li>
            <li class="wbExpo_format_txt"><input type="radio" field="export_format" style="border:none" name="wbExpo_format" value="txt" id="wbExpo_format_txt"/>
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

		for i, col := range colNames {
			buf.WriteString("<li><input type=\"checkbox\" style=\"border:none\" checked=\"checked\" field=\"export_fields[")
			buf.WriteString(strconv.Itoa(i + 1))
			buf.WriteString("]\" id=\"export_field_")
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

	buf.WriteString(`<iframe id="export_frame" style="display:none"></iframe>
        <div style="clear:both"></div></div><input type="button" class="gra-btn btn-export" onclick="wbExpo.doExport('`)

	buf.WriteString(portal.PortalKey)
	buf.WriteString(`')" value=" 导出 "/>`)

	return buf.String()

}
