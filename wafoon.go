package main

import (
    "os"
    "fmt"
    "flag"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/service/waf"

    "github.com/olekukonko/tablewriter"
)

const AppVersion = "0.0.1"

var (
    argProfile = flag.String("profile", "", "Profile 名を指定.")
    argRegion = flag.String("region", "ap-northeast-1", "Region 名を指定.")
    argEndpoint = flag.String("endpoint", "", "AWS API のエンドポイントを指定.")
    argVersion = flag.Bool("version", false, "バージョンを出力.")
    argAclid = flag.String("aclid", "", "Web ACL ID を指定.")
    argAllow = flag.Bool("allow", false, "Default Action を Allow に変更.")
    argBlock = flag.Bool("block", false, "Default Action を Block に変更")
)

func awsWafClient(profile string, region string) *waf.WAF {
    var config aws.Config
    if profile != "" {
        creds := credentials.NewSharedCredentials("", profile)
        config = aws.Config{Region: aws.String(region), Credentials: creds, Endpoint: aws.String(*argEndpoint)}
    } else {
        config = aws.Config{Region: aws.String(region), Endpoint: aws.String(*argEndpoint)}
    }
    sess := session.New(&config)
    wafClient := waf.New(sess)
    return wafClient
}

func getChangeToken(wafClient *waf.WAF) (changeToken string) {
    input := &waf.GetChangeTokenInput{}
    res, err := wafClient.GetChangeToken(input)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
    changeToken = *res.ChangeToken

    return changeToken
}

func updateWebAcl(wafClient *waf.WAF, webAclId string, changeToken string ,defaultActionType string) {
    input := &waf.UpdateWebACLInput{
        ChangeToken: aws.String(changeToken),
        DefaultAction: &waf.WafAction{
            Type: aws.String(defaultActionType),
        },
        WebACLId: aws.String(webAclId),
    }

    fmt.Print("処理を続行しますか? (y/n): ")
    var stdin string
    fmt.Scan(&stdin)
    switch stdin {
    case "y", "Y":
        fmt.Println("処理を続行します.")
    case "n", "N":
        fmt.Println("処理を停止します.")
        os.Exit(0)
    default:
        fmt.Println("処理を停止します.")
        os.Exit(0)
    }

    _, err := wafClient.UpdateWebACL(input)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    } else {
        fmt.Println("デフォルトアクションを変更しました.")
        listWebAcl(wafClient)
    }
}

func outputTbl(data [][]string) {
    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"Name", "WebACLId", "DefaultAction"})
    for _, value := range data {
        table.Append(value)
    }
    table.Render()
}

func getWebAclDefaultAction(wafClient *waf.WAF, webAclId string) (defaultActionType string) {
    input := &waf.GetWebACLInput{
        WebACLId: aws.String(webAclId),
    }
    res, err := wafClient.GetWebACL(input)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    defaultActionType = *res.WebACL.DefaultAction.Type
    return defaultActionType
}

func listWebAcl(wafClient *waf.WAF) {
    input := &waf.ListWebACLsInput{
        Limit: aws.Int64(1),
    }

    allWebAcl := [][]string{}
    for {
        res, err := wafClient.ListWebACLs(input)
        if err != nil {
            fmt.Println(err.Error())
            os.Exit(1)
        }

        for _, r := range res.WebACLs {
            defaultActionType := getWebAclDefaultAction(wafClient, *r.WebACLId)
            webAcl := []string{
                *r.Name,
                *r.WebACLId,
                defaultActionType,
            }
            allWebAcl = append(allWebAcl, webAcl)
        }

        if res.NextMarker == nil {
            break
        }
        input.SetNextMarker(*res.NextMarker)
        continue
    }

    outputTbl(allWebAcl)
}

func main() {
    flag.Parse()

    if *argVersion {
        fmt.Println(AppVersion)
        os.Exit(0)
    }

    wafClient := awsWafClient(*argProfile, *argRegion)

    if *argAllow {
        if *argAclid != "" {
            changeToken := getChangeToken(wafClient)
            updateWebAcl(wafClient, *argAclid, changeToken, "ALLOW")
        } else {
            fmt.Println("Web ACL ID を入力してください.")
            os.Exit(1)
        }
    } else if *argBlock {
        if *argAclid != "" {
            changeToken := getChangeToken(wafClient)
            updateWebAcl(wafClient, *argAclid, changeToken, "BLOCK")
        } else {
            fmt.Println("Web ACL ID を入力してください.")
            os.Exit(1)
        }
    } else {
        listWebAcl(wafClient)
    }
}
