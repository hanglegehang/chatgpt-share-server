package config

import (
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gview"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

var (
	PORT      = 8001
	CHATPROXY = "https://demo.xyhelper.cn"
	AUTHKEY   = "xyhelper"
	ArkoseUrl = "/v2/"

	AssetPrefix  = "https://oaistatic-cdn.closeai.biz"
	BuildId      = "-wRE4Obkm_QOW7PLn1x21"
	CacheBuildId = "-wRE4Obkm_QOW7PLn1x21"
	Script       = "https://cdn.oaistatic.com/_next/static/chunks/2565-263427db2ed7a61a.js?dpl=37f91bfd782f6b4fb81dd5cd885a42d5d31cc4a3"
	Dpl          = "dpl=37f91bfd782f6b4fb81dd5cd885a42d5d31cc4a3"
	envScriptTpl = `
	<script src="/jquery.min.js"></script>
	<script src="/list.js"></script>

	<script>
	window.__arkoseUrl="{{.ArkoseUrl}}";
	window.__assetPrefix="{{.AssetPrefix}}";
	window.__script="{{.Script}}";
	window.__dpl="{{.Dpl}}";
	</script>
	`
	OauthUrl              = ""
	AuditLimitUrl         = ""
	APIAUTH               = ""
	DISALLOW_ROAM         = false // 是否禁止漫游
	FILESERVER            = "https://files.closeai.biz"
	ConversationNotifyUrl = ""
	// Generator *badge.Generator
)

func init() {
	ctx := gctx.GetInitCtx()
	chatproxy := g.Cfg().MustGetWithEnv(ctx, "CHATPROXY").String()
	if chatproxy != "" {
		CHATPROXY = chatproxy
	}
	g.Log().Info(ctx, "CHATPROXY:", CHATPROXY)
	authkey := g.Cfg().MustGetWithEnv(ctx, "AUTHKEY").String()
	if authkey != "" {
		AUTHKEY = authkey
	}
	g.Log().Info(ctx, "AUTHKEY:", AUTHKEY)
	arkoseUrl := g.Cfg().MustGetWithEnv(ctx, "ARKOSE_URL")
	if !arkoseUrl.IsEmpty() {
		ArkoseUrl = arkoseUrl.String()
	}
	assetPrefix := g.Cfg().MustGetWithEnv(ctx, "ASSET_PREFIX").String()
	if assetPrefix != "" {
		AssetPrefix = assetPrefix
	}
	g.Log().Info(ctx, "ASSET_PREFIX:", AssetPrefix)
	cacheBuildId := CheckVersion(ctx, AssetPrefix)
	if cacheBuildId != "" {
		CacheBuildId = cacheBuildId
	}
	g.Log().Info(ctx, "CacheBuildId:", CacheBuildId)
	build, script, dpl := CheckNewVersion(ctx)
	if build != "" {
		BuildId = build
	}
	if script != "" {
		Script = script
	}
	if dpl != "" {
		Dpl = dpl
	}
	g.Log().Info(ctx, "BuildId:", BuildId, "Script:", Script, "Dpl:", Dpl)
	port := g.Cfg().MustGetWithEnv(ctx, "PORT").Int()
	if port != 0 {
		PORT = port
	}
	g.Log().Info(ctx, "PORT:", PORT)
	s := g.Server()
	s.SetPort(PORT)
	s.SetServerRoot("resource/public")
	disallowRoam := g.Cfg().MustGetWithEnv(ctx, "DISALLOW_ROAM").Bool()
	if disallowRoam {
		DISALLOW_ROAM = disallowRoam
	}

	oauthUrl := g.Cfg().MustGetWithEnv(ctx, "OAUTH_URL").String()
	if oauthUrl != "" {
		OauthUrl = oauthUrl
	} else {
		OauthUrl = "http://127.0.0.1:" + gconv.String(PORT) + "/auth/oauth"
	}
	g.Log().Info(ctx, "OAUTH_URL:", OauthUrl)
	auditLimitUrl := g.Cfg().MustGetWithEnv(ctx, "AUDIT_LIMIT_URL").String()
	if auditLimitUrl != "" {
		AuditLimitUrl = auditLimitUrl
	}
	g.Log().Info(ctx, "AUDIT_LIMIT_URL:", AuditLimitUrl)
	conversationNotifyUrl := g.Cfg().MustGetWithEnv(ctx, "ConversationNotifyUrl").String()
	if conversationNotifyUrl != "" {
		ConversationNotifyUrl = conversationNotifyUrl
	}
	g.Log().Info(ctx, "ConversationNotifyUrl:", ConversationNotifyUrl)
	apiAuth := g.Cfg().MustGetWithEnv(ctx, "APIAUTH").String()
	if apiAuth != "" {
		APIAUTH = apiAuth
	}
	g.Log().Info(ctx, "APIAUTH:", APIAUTH)
	fileserver := g.Cfg().MustGetWithEnv(ctx, "FILESERVER").String()
	if fileserver != "" {
		FILESERVER = fileserver
	}
	g.Log().Info(ctx, "FILESERVER:", FILESERVER)

	// 每10分钟检查一次版本
	go func() {
		for {
			build, script, dpl := CheckNewVersion(ctx)
			if build != "" {
				BuildId = build
			}
			if script != "" {
				Script = script
			}
			if dpl != "" {
				Dpl = dpl
			}
			g.Log().Info(ctx, "BuildId:", BuildId, "Script:", Script, "Dpl:", Dpl)
			cacheBuildId := CheckVersion(ctx, AssetPrefix)
			if cacheBuildId != "" {
				CacheBuildId = cacheBuildId
			}
			g.Log().Info(ctx, "CacheBuildId:", CacheBuildId)
			g.Log().Info(ctx, "CheckNewVersion:", BuildId, CacheBuildId)
			time.Sleep(10 * time.Minute)
		}

	}()
}

func GetEnvScript(ctx g.Ctx) string {
	script, err := gview.ParseContent(ctx, envScriptTpl, g.Map{
		"ArkoseUrl":   ArkoseUrl,
		"AssetPrefix": AssetPrefix,
		"Script":      Script,
		"Dpl":         Dpl,
	})
	if err != nil {
		g.Log().Error(ctx, "GetEnvScript Error: ", err)
		return ""
	}
	return script
}

// 检查是否有新版本
func CheckNewVersion(ctx g.Ctx) (buildId, script, dpl string) {
	// 读取 https://tcr9i-dev.closeai.biz/ping
	respVar := g.Client().GetVar(ctx, "https://publicapi.closeai.biz/chatgpt/info")
	respJson := gjson.New(respVar)
	respJson.Dump()
	buildId = gjson.New(respVar).Get("buildId").String()
	script = gjson.New(respVar).Get("script").String()
	dpl = gjson.New(respVar).Get("dpl").String()
	g.Log().Info(ctx, "Check tcr9i buildId: ", buildId)
	return
}

// 检查版本号并同步资源
func CheckVersion(ctx g.Ctx, assetPrefix string) (CacheBuildId string) {
	gclient := g.Client()
	// 读取 assetPrefix/version
	versionVar := gclient.GetVar(ctx, assetPrefix+"/version.json")
	CacheBuildId = gjson.New(versionVar).Get("cacheBuildId").String()
	g.Log().Infof(ctx, "Get config From %s ,CacheBuildId: %s", AssetPrefix, CacheBuildId)
	if CacheBuildId == "" {
		return ""
	}
	// 读取buildDate目录索引
	indexUrl := assetPrefix + "/template/" + CacheBuildId + "/index.txt"
	g.Log().Info(ctx, "Get config From ", indexUrl)
	buildDateVar := gclient.GetVar(ctx, indexUrl).String()
	if buildDateVar == "" {
		return ""
	}
	// 按回车分割
	buildDateList := gstr.Split(buildDateVar, "\n")
	// g.Dump(buildDateList)
	// 遍历目录索引 如果没有就下载
	for _, v := range buildDateList {
		if v == "" {
			continue
		}
		// 检查文件是否存在
		if !gfile.Exists("./resource/template/" + CacheBuildId + "/" + v) {
			g.Log().Infof(ctx, "Download %s", v)
			// 下载文件
			res, err := gclient.Get(ctx, assetPrefix+"/template/"+CacheBuildId+"/"+v)
			if err != nil {
				g.Log().Error(ctx, "Download  Error: ", v, err)
				return ""
			}
			defer res.Close()
			if res.StatusCode != 200 {
				g.Log().Error(ctx, "Download  Error: ", v, res.StatusCode)
				return ""
			}
			// 写入文件
			err = gfile.PutBytes("./resource/template/"+CacheBuildId+"/"+v, res.ReadAll())
			if err != nil {
				g.Log().Error(ctx, "Download  Error: ", v, err)
				return ""
			}

		}
	}

	return
}
