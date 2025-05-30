package i18nc

import "github.com/nicksnyder/go-i18n/v2/i18n"

func (ic *I18nCustom) AddConfigs() {
	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "ReplyNum",
		One:   "{{.Count}} reply",
		Other: "{{.Count}} replies",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "AddNew",
		Other: "New",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Login",
		Other: "Login",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Register",
		Other: "Register",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Logout",
		Other: "Logout",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Settings",
		One:   "Setting",
		Other: "Settings",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Permission",
		One:   "Permission",
		Other: "Permissions",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Role",
		One:   "Role",
		Other: "Roles",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "User",
		One:   "User",
		Other: "Users",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Activity",
		One:   "Activity",
		Other: "Activities",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Best",
		Other: "Best",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Latest",
		Other: "Latest",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Hot",
		Other: "Hot",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "PublishInfo",
		Other: "By {{.Username}} ",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "VoteScore",
		Other: "vote score {{.Score}}",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Discuss",
		Other: "discuss",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Weight",
		Other: "weight {{.Weight}}",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Participate",
		One:   "{{.ParticipateNum}} participate",
		Other: "{{.ParticipateNum}} participates",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "UISaveSuccess",
		Other: "UI settings successfully saved",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Account",
		Other: "Account",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "UI",
		Other: "UI",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Username",
		Other: "Username",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Introduction",
		Other: "Introduction",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Language",
		Other: "Language",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Theme",
		Other: "Theme",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "ThemeLight",
		Other: "Light",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "ThemeDark",
		Other: "Dark",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "ThemeSystem",
		Other: "OS Default",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "EnableJavaScriptTip",
		Other: "Must enable JavaScript",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "PageLayout",
		Other: "Page Layout",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "PageLayoutFull",
		Other: "Full",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "PageLayoutCentered",
		Other: "Centered",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "AccountSaveSuccess",
		Other: "Account settings successfully saved",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "AccountCreateSuccess",
		Other: "Account created successfully",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Type",
		Other: "Type",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Action",
		Other: "Action",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "All",
		Other: "All",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Manage",
		Other: "Manage",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Re",
		Other: "Re",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Link",
		Other: "Link",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Anchor",
		Other: "Anchor",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Deleted",
		Other: "Deleted",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "ConfirmDelete",
		Other: "Confirm to delete",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "ReactTip",
		Other: "React to content",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "AddContent",
		Other: "Add content",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Content",
		Other: "Content",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Title",
		Other: "Title",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Modified",
		Other: "Modified",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Email",
		Other: "Email",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Password",
		Other: "Password",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "PasswordFormatTip",
		Other: "Password must be at least {{.LeastLen}} characters long and contain a combination of numbers, letters, and special characters.",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "UsernameFormatTip",
		Other: "Default to using the username from the email. Username can only consist of numbers, letters, and characters _.-, and cannot begin or end with a symbol.",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "LoginTip",
		Other: "Already have an account? Please {{.LoginLink}} directly.",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Or",
		Other: "{{.A}} or {{.B}}",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "RegisterTipHead",
		Other: "Not registered yet? You can ",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "RegisterTip",
		Other: "Create a new account",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "GoTo",
		Other: "Go to ",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "GoBack",
		Other: "Go back",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "GoHome",
		Other: "Go home",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Author",
		Other: "Author",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Article",
		One:   "Article",
		Other: "Articles",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "ArticleTitle",
		Other: "Article title",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "ArticleContent",
		Other: "Article content",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "UserManage",
		Other: "{{local \"User\"}} {{local \"Manage\"}}",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "UpdateRole",
		Other: "Update {{local \"Role\"}}",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Reply",
		One:   "Reply",
		Other: "Replies",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Saved",
		Other: "Saved",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "JoinAt",
		Other: "Joined At",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "NoData",
		Other: "No data",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "PublishSuccess",
		Other: "Content published successfully",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "DeleteSuccess",
		Other: "Content deleted successfully",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "EditContent",
		Other: "Edit content",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "UserList",
		Other: "User List",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "List",
		Other: "{{.Name}} List",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "AddItem",
		Other: "Add {{.Name}}",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "EditItem",
		Other: "Edit {{.Name}}",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Oldest",
		Other: "Oldest",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "URL",
		Other: "URL",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Source",
		Other: "Source",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "FormOptional",
		Other: "Optional",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "FormRequired",
		Other: "Required",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "ArticleTitleTip",
		Other: "Up to {{.Num}} characters, please summarize content concisely without clickbait titles. Irrelevant content will be removed.",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "ArticleURLTip",
		Other: "Please provide direct links, avoid using redirected URLs. Whenever possible, provide primary sources.",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "ArticleContentTip",
		Other: "Up to {{.Num}} characters.",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "FontSize",
		Other: "Font Size",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "FontExtremSmall",
		Other: "Extrem Small",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "FontSmall",
		Other: "Small",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "FontRegular",
		Other: "Regular",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "FontLarge",
		Other: "Large",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "FontExtremLarge",
		Other: "Extrem Large",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "FontCustom",
		Other: "Custom",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Upvote",
		Other: "Upvote",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Downvote",
		Other: "Downvote",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "CancelVote",
		Other: "Cancel the vote",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "SkipToContent",
		Other: "Skip to content",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Subscribed",
		Other: "Subscribed",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "MessageUnread",
		Other: "Unread",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "MessageRead",
		Other: "Read",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Message",
		Other: "Message",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "NewReply",
		Other: "New reply on {{.ArticleTitle}}",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "EmailVerify",
		Other: "Email Verification",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "VerificationEmailTip",
		Other: "The verification code has been sent to the email: {{.Email}}. It is valid for {{.Duration}} minutes. Please enter the code to complete the registration.",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "VerificationCode",
		Other: "Verification code",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "ResendVerification",
		Other: "Resend the verification code to the email.",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "VerificationExpired",
		Other: "The verification code has expired.",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "VerificationIncorrect",
		Other: "The verification code is incorrect.",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "SubmitContentTip",
		Other: "Due to the content being published on the internet, please refrain from including personal privacy information in the post title and content. All private data will be removed.",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID: "VerificationMailTpl",
		Other: `<html>
<body>
<p>You are registering on {{.DomainName}}, here's the verfication code:</p>
<p><large><b>{{.Code}}</b></large></p>
<p>Valid for {{.Minutes}} minutes.</p>
<hr>
<p style="color:#666">{{.DomainName}}</p>
</body>
</html>`,
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID: "VerificationResetPassMailTpl",
		Other: `<html>
<body>
<p>You are resetting the password on {{.DomainName}}, here's the verfication code:</p>
<p><large><b>{{.Code}}</b></large></p>
<p>Valid for {{.Minutes}} minutes.</p>
<hr>
<p style="color:#666">{{.DomainName}}</p>
</body>
</html>`,
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "VerificationMailTitle",
		Other: "Verification code for registration",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "VerificationResetPassMailTitle",
		Other: "Verification code for resetting password",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "AcountExistsTip",
		Other: "This account has already been registered on this platform. Please log in using an alternative method.",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "OAuthLoginTip",
		Other: "or log in using the following platform",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "RetrievePassword",
		Other: "Retrieve password",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "PassResetSuccess",
		Other: "Password reset successful",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "RetrievePassTip",
		Other: "Please enter the email associated with your account.",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "ResetPassTip",
		Other: "If a matching account is detected, the verification code will be sent to the email: {{.Email}}, valid for {{.Duration}} minute. Please enter the new password and the verification code to complete the password reset.",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "NewPassword",
		Other: "New password",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "ConfirmNewPassword",
		Other: "Confirm new password",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Reason",
		Other: "Reason",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Emoji",
		Other: "Emoji",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "ShowItem",
		Other: "Show {{.Name}}",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "BrandName",
		Other: "DizKaz",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Category",
		One:   "Category",
		Other: "Categories",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "ResetPassword",
		Other: "Reset Password",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "NewArticleInCategory",
		Other: "{{.AuthorName}} publised new article {{.ArticleTitle}} under {{.CategoryName}}",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Version",
		Other: "Version",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "EditHistoryTitle",
		Other: "Edit history of {{.Title}}",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "EditBy",
		Other: "Edit by {{.Name}} ",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "About",
		Other: "About",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "TimeOrder",
		Other: "Chronological",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "RepliesLayout",
		Other: "Replies Layout",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "RepliesLayoutTree",
		Other: "Tree",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "RepliesLayoutTile",
		Other: "Tile",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Pin",
		Other: "Pin",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "PinExpireAt",
		Other: "Pin expires at {{.Time}}",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "PinExpireTime",
		Other: "Pin expires time",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Lock",
		Other: "Lock",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Locked",
		Other: "Locked",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "HideChanges",
		Other: "Hide changes from edit history",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "EditHistory",
		Other: "Edit history",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "EditHistoryHidden",
		Other: "The infraction content has been removed",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Trash",
		Other: "Trash",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Keyword",
		One:   "Keyword",
		Other: "Keywords",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "ShareTip",
		Other: "Please copy the above link and share it",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Share",
		Other: "Share",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "MainlandChina",
		Other: "Mainland China",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "UnitedStates",
		Other: "United States",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "India",
		Other: "India",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "BlockRegionsTip",
		Other: "Please select regions to block",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "BlockedRegions",
		Other: "Blocked Regions",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "ArticleListDefaultSort",
		Other: "Article List Default Sort Type",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "ReplyListDefaultSort",
		Other: "Reply List Default Sort Type",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Reputation",
		Other: "Reputation",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Guide",
		Other: "Guide",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "PleaseSelect",
		Other: "-- select an option --",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "AlreadyBan",
		Other: "Already banned",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "AlreadyUnban",
		Other: "Already unbanned",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "ConfirmBan",
		Other: "Confirm to ban {{.Name}}?",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "ConfirmUnban",
		Other: "Confirm to unban {{.Name}}?",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "UnbanTime",
		Other: "Unban time",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "UnitDay",
		One:   "{{.Count}} day",
		Other: "{{.Count}} days",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Forever",
		Other: "Forever",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "BannedDuration",
		Other: "Banned duration",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "BannedTimes",
		Other: "Banned times",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "BannedStatusTip",
		Other: "This account is banned for {{.CountDays}}",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "BannedForeverTip",
		Other: "This account is banned forever",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "UnbanSuccessTip",
		Other: "Unbanned successfully",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "AnalysisReport",
		Other: "Analysis Report",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Voted",
		Other: "Voted",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Matrix",
		Other: "Matrix",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "AuthFrom",
		Other: "Auth From",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "CharCount",
		Zero:  "No Content",
		Other: "Content Character Count {{.Count}}",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "Refresh",
		Other: "Refresh",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "SearchSite",
		Other: "Search",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "ResponseTime",
		Other: "Response time",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "HTMLResponseTime",
		Other: "HTML rendering time",
	})

	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "AboutSiteTitle",
		Other: "About {{.SiteName}}",
	})


	ic.AddLocalizeConfig(&i18n.Message{
		ID:    "AboutSiteContent",
		Other: "<p>{{.SiteName}} is an online community where content is ranked by user votes. The inspiration for this community came from Ruan Yifeng's series of articles on <a href=\"https://www.ruanyifeng.com/blog/2012/02/ranking_algorithm_hacker_news.html\" target=\"_blank\">ranking algorithms based on user votes</a>, combined with the developer's own interest in the topic. If you enjoy sharing and discovering quality content online and participating in thoughtful discussions about it, please feel free to <a href=\"/register\">join</a> us.</p>",
	})
}
