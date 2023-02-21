package constants

// TODO: this is just a demo, need plenish
const HotfixPendingApprovalMsg = "There is a new Hotfix creation ticket!\n" +
	"**Name**: %[1]s\n" +
	"**Customer**: %[2]s\n" +
	"**Related Repo**: %[3]s\n" +
	"**Related Issues**: %[4]s\n" +
	"**Related PullRequests**: %[5]s\n" +
	"You can click below link to check more details.\n"

const HotfixApproveMsg = "Your aplication for new hotfix %[1]s has been approved!\n" +
	"You can click below link to start triage.\n"

const HotfixDenyMsg = "Your aplication for new hotfix %[1]s has been denied.\n" +
	"Please contact RM for more info.\n"
