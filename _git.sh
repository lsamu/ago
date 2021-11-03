read -p "请输入本次提交的描述：" remark
git config user.name "lauxinyi"
git config user.email "lauxinyi@gmail.com"
git pull
git status
git add -A
git commit -m "$remark"
git push 