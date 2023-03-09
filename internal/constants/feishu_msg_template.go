package constants

// TODO: this is just a demo, need plenish
const HotfixPendingApprovalMsg = "**Info**\n" +
	"**Name**: %[1]s     **Customer**: %[2]s \n" +
	"**Submitor**: %[3]s     **Oncall**: <a href='%[4]s'>%[5]s</a> \n" +
	"**Based Version**: %[6]s   \n" +
	"**Related Repos**: %[7]s   \n" +
	"You can check the link bellow for more information.\n"

const HotfixApproveMsg = "Your aplication for new hotfix %[1]s has been approved!\n" +
	"You can click below link to start triage.\n"

const HotfixDenyMsg = "Your aplication for new hotfix %[1]s has been denied.\n" +
	"Please contact RM for more info.\n"

const TiReleaseUrl = "https://tirelease.pingcap.net"
