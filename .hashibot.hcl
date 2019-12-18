behavior "regexp_issue_labeler" "bug_label" {
    regexp = " \\[issue-type:bug-report\\] "
    labels = ["bug"]
}

behavior "regexp_issue_labeler" "enhancement_label" {
    regexp = " \\[issue-type:enhancement\\] "
    labels = ["enhancement"]
}

poll "stale_issue_closer" "closer" {
    schedule = "0 50 12 * * *"
    labels = ["stale"]
    no_reply_in_last = "2160h" # 90 days
    max_issues = 500
    sleep_between_issues = "5s"
    message = <<-EOF
    I'm going to close this issue due to inactivity (_90 days_ without response ⏳ ). This helps our maintainers find and focus on the active issues.

    If you feel I made an error 🤖 🙉  , please reach out to my human friends 👉  hashibot-feedback@hashicorp.com. Thanks!
    EOF
}

poll "closed_issue_locker" "locker" {
    schedule = "0 50 13 * * *"
    closed_for = "720h" # 30 days
    max_issues = 500
    sleep_between_issues = "5s"
    message = <<-EOF
    I'm going to lock this issue because it has been closed for _30 days_ ⏳. This helps our maintainers find and focus on the active issues.

    If you feel this issue should be reopened, we encourage creating a new issue linking back to this one for added context. If you feel I made an error 🤖 🙉  , please reach out to my human friends 👉  hashibot-feedback@hashicorp.com. Thanks!
    EOF
}

behavior "assign_random_reviewer" "random" {
  reviewers            = [
    "paddycarver",
    "danawillow",
    "megan07",
    "rileykarson",
    "tysen",
    "ndmckinley",
    "slevenick",
    "emilymye",
    "chrisst",
  ]
  only_non_maintainers = true
}
