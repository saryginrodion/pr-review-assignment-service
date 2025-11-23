from locust import HttpUser, task, between
import random
import string


def random_id(prefix):
    return prefix + "".join(random.choices(string.ascii_lowercase + string.digits, k=6))


class PRReviewerLoadTest(HttpUser):
    wait_time = between(0.1, 0.3)

    def on_start(self):
        self.team_name = random_id("team_")

        self.users = [{"user_id": random_id("u_"), "username": f"User{i}", "is_active": True} for i in range(20)]

        self.client.post("/team/add", json={"team_name": self.team_name, "members": self.users})

        self.prs = []

    def get_active_candidates(self, exclude):
        exclude_set = set(exclude)
        return [u["user_id"] for u in self.users if u["is_active"] and u["user_id"] not in exclude_set]

    @task(3)
    def create_pr(self):
        pr_id = random_id("pr_")
        author = random.choice(self.users)["user_id"]

        res = self.client.post("/pullRequest/create", json={"pull_request_id": pr_id, "pull_request_name": "LoadTest PR", "author_id": author})

        if res.status_code == 201:
            pr = res.json()["pr"]
            self.prs.append(pr)

    @task(2)
    def merge_pr(self):
        open_prs = [p for p in self.prs if p["status"] == "OPEN"]
        if not open_prs:
            return

        pr = random.choice(open_prs)

        res = self.client.post("/pullRequest/merge", json={"pull_request_id": pr["pull_request_id"]})

        if res.status_code == 200:
            updated = res.json()["pr"]
            pr["status"] = updated["status"]
            pr["assigned_reviewers"] = updated["assigned_reviewers"]
            pr["mergedAt"] = updated.get("mergedAt")

    @task(2)
    def reassign_valid(self):
        open_prs = [p for p in self.prs if p["status"] == "OPEN"]
        if not open_prs:
            return

        pr = random.choice(open_prs)

        if not pr["assigned_reviewers"]:
            return

        old_user = random.choice(pr["assigned_reviewers"])

        candidates = self.get_active_candidates(exclude=pr["assigned_reviewers"] + [pr["author_id"]])
        if not candidates:
            return

        res = self.client.post("/pullRequest/reassign", json={"pull_request_id": pr["pull_request_id"], "old_reviewer_id": old_user})

        if res.status_code == 200:
            body = res.json()
            updated = body["pr"]

            pr["assigned_reviewers"] = updated["assigned_reviewers"]

    @task(1)
    def get_review(self):
        uid = random.choice(self.users)["user_id"]
        self.client.get("/users/getReview", params={"user_id": uid})

    @task(1)
    def flip_user_active(self):
        user = random.choice(self.users)
        user["is_active"] = not user["is_active"]

        self.client.post("/users/setIsActive", json={"user_id": user["user_id"], "is_active": user["is_active"]})
