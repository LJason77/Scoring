webpackJsonp([1],{NHnr:function(t,a,e){"use strict";Object.defineProperty(a,"__esModule",{value:!0});var i=e("7+uW"),l={render:function(){var t=this.$createElement,a=this._self._c||t;return a("div",{attrs:{id:"app"}},[a("router-view")],1)},staticRenderFns:[]};var s=e("VU/8")({name:"App"},l,!1,function(t){e("UcgC")},null,null).exports,n=e("/ocq"),r=e("Dd8w"),o=e.n(r),c=e("NYxO"),d={name:"",data:function(){return{timer:Object,isgetData:!0,input:"",tableData:[],detailNum:0,detailData:[],CCdetailData:[],apidetailData:[]}},computed:o()({},Object(c.b)([])),methods:o()({},Object(c.a)([]),{renderHeader:function(t,a){a.column;return t("div",{style:"cursor:pointer;",on:{click:this.$refs.teamTable.clearSort}},"队伍编号")},shownDetail:function(t){this.detailData=t,this.detailNum=t.teamNum,this.CCdetailData=[],this.apidetailData=[];for(var a=0;a<t.CCDetails.length;a++)this.CCdetailData.push(t.CCDetails[a].split(":"));for(a=0;a<t.apiDetails.length;a++)this.apidetailData.push(t.apiDetails[a].split(":"))},CCtableRowClassName:function(t){var a=t.rowIndex;return-1!==this.CCdetailData[a][1].indexOf("fail")?"warning-row":-1!==this.CCdetailData[a][1].indexOf("pass")?"success-row":""},apitableRowClassName:function(t){var a=t.rowIndex;return-1!==this.apidetailData[a][1].indexOf("fail")?"warning-row":-1!==this.apidetailData[a][1].indexOf("pass")?"success-row":""},preview:function(t){var a=document.getElementById(this.id).files[0];this.imgDataUrl=this.getObjectURL(a),this.$emit("sendImgUrl",this.imgDataUrl,this.id)},getObjectURL:function(t){var a=null;return void 0!=window.createObjectURL?a=window.createObjectURL(t):void 0!=window.webkitURL?a=window.webkitURL.createObjectURL(t):void 0!=window.URL&&(a=window.URL.createObjectURL(t)),a},clickFun:function(t){var a=this;""!=this.input?this.$axios.get("/invalidPath?path="+this.input).then(function(t){if(2==t.data.code)if(null!=t.data.data){for(var e="",i=0;i<t.data.data.length;i++)e+=t.data.data[i]+"<br/>";a.$alert(e,"缺少路径",{dangerouslyUseHTMLString:!0})}else a.$message({message:"无缺少路径",type:"success"});else a.$message.error(t.data.message)}).catch(function(t){return a.$message.error("项目文件夹不存在")}):this.$message.error("请输入地址")},stopTest:function(){clearInterval(this.timer),this.isgetData=!0},getData:function(){if(""!=this.input){var t=this;t.isgetData=!1,t.timer=setInterval(function(){1!=t.isgetData?t.goTest():clearInterval(t.timer)},1e3)}else this.$message.error("请输入地址")},goTest:function(){var t=this;this.$axios.get("/scoring?path="+this.input).then(function(a){if(1==a.data.code);else if(2==a.data.code){t.tableData=a.data.data;for(var e=0;e<t.tableData.length;e++)t.tableData[e].apipassNum=Math.ceil(t.tableData[e].apiPassingRate*t.tableData[e].apiDetails.length),t.tableData[e].CCpassNum=Math.ceil(t.tableData[e].CCPassingRate*t.tableData[e].CCDetails.length),t.tableData[e].apiPassingRate=100*t.tableData[e].apiPassingRate.toFixed(4)+"%",t.tableData[e].CCPassingRate=100*t.tableData[e].CCPassingRate.toFixed(4)+"%",t.tableData[e].total=t.tableData[e].apiScore+t.tableData[e].CCScore,t.tableData[e].teamNum=e+1,t.apidetailData=[],t.CCdetailData=[],t.isgetData=!0,t.detailData=[],t.detailNum=0;t.$message({message:"测试完成",type:"success"})}else 3==a.data.code&&(t.$message.error(a.data.message),t.isgetData=!0)}).catch(function(a){t.$message.error("没有此文件路径"),t.isgetData=!0})}}),beforeCreate:function(){var t=this;this.$axios.get("/scoring?path=./projects").then(function(a){if(2==a.data.code){t.tableData=a.data.data;for(var e=0;e<t.tableData.length;e++)t.tableData[e].apipassNum=Math.ceil(t.tableData[e].apiPassingRate*t.tableData[e].apiDetails.length),t.tableData[e].CCpassNum=Math.ceil(t.tableData[e].CCPassingRate*t.tableData[e].CCDetails.length),t.tableData[e].apiPassingRate=100*t.tableData[e].apiPassingRate.toFixed(4)+"%",t.tableData[e].CCPassingRate=100*t.tableData[e].CCPassingRate.toFixed(4)+"%",t.tableData[e].total=t.tableData[e].apiScore+t.tableData[e].CCScore,t.tableData[e].teamNum=e+1}}).catch(function(t){console.log(t)})}},u={render:function(){var t=this,a=t.$createElement,e=t._self._c||a;return e("el-row",{staticClass:"center"},[e("el-col",{staticClass:"left",attrs:{span:17}},[e("div",{staticClass:"left-title"},[e("div",{staticClass:"left-title-one"},[e("div",{staticClass:"left-title-one-text1"},[t._v("1输入总地址")]),t._v(" "),e("div",{staticClass:"left-title-one-text2"},[e("el-input",{attrs:{placeholder:"将地址粘贴到此,例如：C:\\Program Files\\Java\\jre1.8.0_201"},model:{value:t.input,callback:function(a){t.input=a},expression:"input"}})],1),t._v(" "),e("el-row",{staticClass:"left-title-one-text-btn"},[e("el-button",{attrs:{type:"primary",plain:""},on:{click:function(a){return t.clickFun()}}},[t._v("确定")])],1)],1),t._v(" "),e("div",{staticClass:"left-title-two"},[e("div",{staticClass:"left-title-one-text-tisp"},[t._v("（包含全部待测队伍项目的文件地址）")])]),t._v(" "),e("div",{staticClass:"left-title-testBtn"},[e("el-row",{staticClass:"left-title-one-text-btn"},[e("el-button",{attrs:{type:"primary",plain:""},on:{click:function(a){return t.getData()}}},[t._v("开始测试")]),t._v(" "),e("el-button",{attrs:{type:"primary",plain:""},on:{click:function(a){return t.stopTest()}}},[t._v("停止")])],1)],1)]),t._v(" "),e("div",{staticClass:"left-table"},[e("el-table",{directives:[{name:"loading",rawName:"v-loading",value:!t.isgetData,expression:"!isgetData"}],ref:"teamTable",staticStyle:{width:"100%"},attrs:{data:t.tableData,"max-height":"600","element-loading-text":"测试中","element-loading-background":"rgba(f, f, f, 0.8)"}},[e("el-table-column",{attrs:{fixed:"left",prop:"teamNum",align:"center",label:"队伍编号","render-header":t.renderHeader}}),t._v(" "),e("el-table-column",{attrs:{align:"center",label:"智能合约",sortable:"","sort-by":"CCScore","sort-orders":["ascending","descending"]}},[e("el-table-column",{attrs:{prop:"CCScore",align:"center",label:"得分",width:"150"}}),t._v(" "),e("el-table-column",{attrs:{prop:"CCPassingRate",align:"center",label:"通过率",width:"150"}})],1),t._v(" "),e("el-table-column",{attrs:{align:"center",label:"应用层",sortable:"","sort-by":"apiScore","sort-orders":["ascending","descending"]}},[e("el-table-column",{attrs:{prop:"apiScore",align:"center",label:"得分"}}),t._v(" "),e("el-table-column",{attrs:{prop:"apiPassingRate",align:"center",label:"通过率"}})],1),t._v(" "),e("el-table-column",{attrs:{prop:"total",align:"center",label:"项目总分",sortable:"","sort-orders":["ascending","descending"]}}),t._v(" "),e("el-table-column",{attrs:{align:"center",label:"详细"},scopedSlots:t._u([{key:"default",fn:function(a){return[e("el-button",{attrs:{type:"text",size:"small"},nativeOn:{click:function(e){return e.preventDefault(),t.shownDetail(a.row)}}},[t._v("详细")])]}}])})],1)],1)]),t._v(" "),e("el-col",{staticClass:"right",attrs:{span:7}},[e("div",{staticClass:"right-header"},[t._v(t._s(0==t.detailNum?"显示队伍的详细测试信息":"当前显示第"+t.detailNum+"队的测试信息"))]),t._v(" "),e("div",{staticClass:"right-center"},[e("div",{staticClass:"right-center-item"},[e("el-table",{attrs:{data:t.CCdetailData,"max-height":"400","header-cell-style":{height:"60px"},"row-class-name":t.CCtableRowClassName}},[e("el-table-column",{attrs:{align:"center",label:"智能合约"},scopedSlots:t._u([{key:"default",fn:function(a){return[t._v(t._s(t.CCdetailData[a.$index][0]))]}}])})],1),t._v(" "),e("div",{staticClass:"right-center-item-foot"},[0!=t.detailNum?e("div",[t._v("\n            测试共检测"+t._s(t.CCdetailData.length)+"项\n            "),e("br"),t._v("\n            检测通过"+t._s(t.detailData.CCpassNum)+"项\n            "),e("br"),t._v("\n            得分为30*（"+t._s(t.detailData.CCpassNum)+"/"+t._s(t.CCdetailData.length)+"）= "+t._s(t.detailData.CCScore)+"\n          ")]):t._e()])],1),t._v(" "),e("div",{staticClass:"right-center-item"},[e("el-table",{attrs:{data:t.apidetailData,"max-height":"400","header-cell-style":{height:"60px"},"row-class-name":t.apitableRowClassName}},[e("el-table-column",{attrs:{align:"center",label:"应用层"},scopedSlots:t._u([{key:"default",fn:function(a){return[t._v(t._s(t.apidetailData[a.$index][0]))]}}])})],1),t._v(" "),e("div",{staticClass:"right-center-item-foot"},[0!=t.detailNum?e("div",[t._v("\n            测试共检测"+t._s(t.apidetailData.length)+"项\n            "),e("br"),t._v("\n            检测通过"+t._s(t.detailData.apipassNum)+"项\n            "),e("br"),t._v("\n            得分为25*（"+t._s(t.detailData.apipassNum)+"/"+t._s(t.apidetailData.length)+"）= "+t._s(t.detailData.apiScore)+"\n          ")]):t._e()])],1)]),t._v(" "),0!=t.detailNum?e("div",{staticClass:"right-foot"},[t._v("队伍"+t._s(t.detailNum)+"，在“智能合约”与“应用层”两项共得分"+t._s(t.detailData.total))]):t._e()])],1)},staticRenderFns:[]};var p=e("VU/8")(d,u,!1,function(t){e("z77L")},"data-v-0e8df487",null).exports;i.default.use(n.a);var g=new n.a({routes:[{path:"/",name:"Home",component:p}]}),b=e("zL8q"),m=e.n(b),h=(e("tvR6"),e("Fd2+")),f=e("mtWM"),C=e.n(f);i.default.prototype.$axios=C.a,i.default.use(h.a),i.default.use(m.a),i.default.config.productionTip=!1,new i.default({el:"#app",router:g,components:{App:s},template:"<App/>"})},UcgC:function(t,a){},tvR6:function(t,a){},z77L:function(t,a){}},["NHnr"]);
//# sourceMappingURL=app.19dd11b97a95b7ea2a5c.js.map