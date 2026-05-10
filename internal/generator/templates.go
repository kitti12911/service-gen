package generator

var commonTemplates = map[string]string{
	".gitignore": `# test coverage
coverage.out

# generated code
gen/

# local config
config.yml
config.yaml

# air
tmp/
`,
	".markdownlint-cli2.jsonc": `{
    "config": {
        "MD013": false
    },
    "globs": ["**/*.{md,markdown}"],
    "ignores": [".github/CODEOWNERS", "CHANGELOG.md"]
}
`,
	".prettierrc.json": `{
    "tabWidth": 4,
    "trailingComma": "none",
    "useTabs": false
}
`,
	".prettierignore": `CHANGELOG.md
.codex/
.cursor/
gen/
tmp/
`,
	".golangci.yml": `version: "2"

run:
    timeout: 5m

linters:
    enable:
        - bodyclose
        - cyclop
        - errcheck
        - errorlint
        - exhaustive
        - gocritic
        - godot
        - gosec
        - govet
        - ineffassign
        - misspell
        - nilerr
        - noctx
        - prealloc
        - revive
        - sqlclosecheck
        - staticcheck
        - unconvert
        - unparam
        - unused
        - wrapcheck

    settings:
        cyclop:
            max-complexity: 15

        exhaustive:
            default-signifies-exhaustive: true

        gocritic:
            enabled-tags:
                - diagnostic
                - performance
                - style
            disabled-checks:
                - hugeParam
                - rangeValCopy

        gosec:
            excludes:
                - G404

        govet:
            enable-all: true
            disable:
                - fieldalignment

        misspell:
            locale: US

        revive:
            rules:
                - name: exported
                  arguments:
                      - disableStutteringCheck
                - name: var-naming
                - name: blank-imports
                - name: context-as-argument
                - name: error-return
                - name: error-strings
                - name: increment-decrement
                - name: range
                - name: receiver-naming
                - name: unused-parameter
                  disabled: true

        unparam:
            check-exported: false

        wrapcheck:
            ignore-sigs:
                - .Errorf(
                - errors.New(
                - errors.Unwrap(
                - errors.Join(
                - .Wrap(
                - .Wrapf(
                - .WithMessage(
                - .WithMessagef(
                - .WithStack(

    exclusions:
        generated: lax
        presets:
            - comments
            - common-false-positives
            - legacy
            - std-error-handling
        rules:
            - path: _test\.go$
              linters:
                  - errcheck
                  - godot
                  - wrapcheck
            - path: ^internal/(server|config|database)/.*\.go$
              linters:
                  - wrapcheck

formatters:
    enable:
        - gofmt
        - goimports
    settings:
        goimports:
            local-prefixes:
                - {{ .ModulePath }}

issues:
    max-issues-per-linter: 0
    max-same-issues: 0
`,
	".vscode/settings.json": `{
    "[markdown]": {
        "editor.defaultFormatter": "esbenp.prettier-vscode",
        "editor.tabSize": 4,
        "editor.insertSpaces": true,
        "editor.detectIndentation": false
    },
    "[yaml]": {
        "editor.defaultFormatter": "esbenp.prettier-vscode",
        "editor.tabSize": 4,
        "editor.insertSpaces": true,
        "editor.detectIndentation": false
    },
    "[json]": {
        "editor.defaultFormatter": "esbenp.prettier-vscode",
        "editor.tabSize": 4,
        "editor.insertSpaces": true,
        "editor.detectIndentation": false
    },
    "[jsonc]": {
        "editor.defaultFormatter": "esbenp.prettier-vscode",
        "editor.tabSize": 4,
        "editor.insertSpaces": true,
        "editor.detectIndentation": false
    }
}
`,
	".vscode/extensions.json": `{
    "recommendations": [
        "davidanson.vscode-markdownlint",
        "esbenp.prettier-vscode",
        "github.vscode-github-actions",
        "golang.go",
        "ms-azuretools.vscode-containers",
        "ms-azuretools.vscode-docker",
        "ms-vscode.makefile-tools"
    ]
}
`,
	".github/CODEOWNERS": `/.github/ {{ .CodeOwner }}
/.vscode/ {{ .CodeOwner }}
/.golangci.yml {{ .CodeOwner }}
/.markdownlint-cli2.jsonc {{ .CodeOwner }}
/.prettierrc.json {{ .CodeOwner }}
/Makefile {{ .CodeOwner }}
`,
	".github/renovate.json": `{
    "$schema": "https://docs.renovatebot.com/renovate-schema.json",
    "extends": ["config:recommended"],
    "timezone": "Asia/Bangkok",
    "schedule": ["* 0-4 1 * *"],
    "updateNotScheduled": false,
    "enabledManagers": ["gomod", "github-actions"],
    "postUpdateOptions": ["gomodTidy"],
    "reviewersFromCodeOwners": true,
    "assigneesFromCodeOwners": true,
    "assignAutomerge": true
}
`,
	".github/workflows/markdownlint.yml": `name: Markdownlint

on:
    push:
        branches:
            - main
        paths:
            - ".github/workflows/markdownlint.yml"
            - ".markdownlint-cli2.*"
            - ".markdownlint.*"
            - "**/*.md"
            - "**/*.markdown"
    pull_request:
        paths:
            - ".github/workflows/markdownlint.yml"
            - ".markdownlint-cli2.*"
            - ".markdownlint.*"
            - "**/*.md"
            - "**/*.markdown"

permissions:
    contents: read

jobs:
    markdownlint:
        name: Markdownlint
        runs-on: ubuntu-latest
        steps:
            # actions/checkout v6.0.2
            - uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd

            # DavidAnson/markdownlint-cli2-action v23.1.0
            - uses: DavidAnson/markdownlint-cli2-action@6b51ade7a9e4a75a7ad929842dd298a3804ebe8b
`,
	".github/workflows/go-ci.yml": `name: Go CI

on:
    push:
        branches:
            - develop
            - uat
            - main
        paths:
            - ".releaserc.json"
            - ".github/workflows/go-ci.yml"
            - ".github/workflows/release.yml"
            - ".github/workflows/markdownlint.yml"
            - ".golangci.yml"
            - "Dockerfile"
            - "Makefile"
            - "buf.gen.yaml"
            - "config.example.yml"
            - "go.mod"
            - "go.sum"
            - "**/*.go"
            - "**/*.proto"
            - "**/*.sql"
            - "**/*.yml"
            - "**/*.yaml"
            - "!README.md"
    pull_request:
        paths:
            - ".releaserc.json"
            - ".github/workflows/go-ci.yml"
            - ".github/workflows/release.yml"
            - ".github/workflows/markdownlint.yml"
            - ".golangci.yml"
            - "Dockerfile"
            - "Makefile"
            - "buf.gen.yaml"
            - "config.example.yml"
            - "go.mod"
            - "go.sum"
            - "**/*.go"
            - "**/*.proto"
            - "**/*.sql"
            - "**/*.yml"
            - "**/*.yaml"
            - "!README.md"

concurrency:
    group: $__GHA_OPEN__ github.workflow __GHA_CLOSE__-$__GHA_OPEN__ github.event.pull_request.number || github.ref __GHA_CLOSE__
    cancel-in-progress: true

permissions:
    contents: read

jobs:
    generate:
        name: Generate
        runs-on: ubuntu-latest
        container:
            image: zot.kittiaccess.work/kitti12911/image-toolchain@sha256:47355f96a059465947c38aa956da1c4502c11d1e8f53eb2c8b3980ba58983d42
            credentials:
                username: $__GHA_OPEN__ secrets.ZOT_USERNAME __GHA_CLOSE__
                password: $__GHA_OPEN__ secrets.ZOT_TOKEN __GHA_CLOSE__
        steps:
            # actions/checkout v6.0.2
            - uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd

            - name: Generate code
              run: |
                  rm -rf gen/grpc
                  buf generate

            # actions/upload-artifact v7.0.1
            - uses: actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a
              with:
                  name: generated-code
                  path: gen/

    lint:
        name: Lint
        runs-on: ubuntu-latest
        needs:
            - generate
        steps:
            # actions/checkout v6.0.2
            - uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd

            # actions/download-artifact v8.0.1
            - uses: actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c
              with:
                  name: generated-code
                  path: .

            # actions/setup-go v6.4.0
            - uses: actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c
              with:
                  go-version-file: go.mod
                  cache: false

            - name: Go vet
              run: go vet ./...

            # golangci/golangci-lint-action v9.2.0
            - uses: golangci/golangci-lint-action@1e7e51e771db61008b38414a730f564565cf7c20
              with:
                  version: v2.12.1
                  args: --timeout=5m

    test:
        name: Test
        runs-on: ubuntu-latest
        needs:
            - generate
        steps:
            # actions/checkout v6.0.2
            - uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd

            # actions/download-artifact v8.0.1
            - uses: actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c
              with:
                  name: generated-code
                  path: .

            # actions/setup-go v6.4.0
            - uses: actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c
              with:
                  go-version-file: go.mod
                  cache: false

            - name: Test with race detector and coverage
              run: |
                  go list -f '__TPL_OPEN__.Dir__TPL_CLOSE__' ./... \
                      | grep -v '/gen/' \
                      | sed "s#^${PWD}#.#" \
                      | xargs go test -race -coverprofile=coverage.out -covermode=atomic

            - name: Remove generated code from coverage
              run: |
                  awk 'NR == 1 || ($1 !~ /\/gen\// && $1 !~ /(_gen|_generated)\.go:/)' coverage.out > coverage.filtered.out
                  mv coverage.filtered.out coverage.out

            - name: Show coverage
              run: |
                  coverage_total="$(go tool cover -func=coverage.out | tail -n 1)"
                  echo "${coverage_total}"
                  {
                      echo "## Go coverage"
                      echo
                      echo '` + "```text" + `'
                      echo "${coverage_total}"
                      echo '` + "```" + `'
                  } >> "${GITHUB_STEP_SUMMARY}"

    toolchain-security:
        name: Toolchain Security
        runs-on: ubuntu-latest
        needs:
            - generate
        container:
            image: zot.kittiaccess.work/kitti12911/security-toolchain@sha256:e05509dad7e83d8f54bbf648584ab89e6aabd6d07adff200bea2180bba702e49
            credentials:
                username: $__GHA_OPEN__ secrets.ZOT_USERNAME __GHA_CLOSE__
                password: $__GHA_OPEN__ secrets.ZOT_TOKEN __GHA_CLOSE__
        steps:
            # actions/checkout v6.0.2
            - uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd
              with:
                  fetch-depth: 0

            # actions/download-artifact v8.0.1
            - uses: actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c
              with:
                  name: generated-code
                  path: .

            - name: govulncheck
              run: govulncheck ./...

            - name: Semgrep
              run: semgrep scan --config=p/golang --config=p/secrets --error

    security:
        name: Security
        runs-on: ubuntu-latest
        steps:
            # actions/checkout v6.0.2
            - uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd
              with:
                  fetch-depth: 0

            # aquasecurity/trivy-action v0.36.0
            - uses: aquasecurity/trivy-action@a9c7b0f06e461e9d4b4d1711f154ee024b8d7ab8
              with:
                  scan-type: fs
                  scan-ref: .
                  scanners: vuln,secret,misconfig
                  exit-code: 1
                  severity: CRITICAL,HIGH
                  ignore-unfixed: true

            # gitleaks/gitleaks-action v2.3.9
            - uses: gitleaks/gitleaks-action@ff98106e4c7b2bc287b24eaf42907196329070c7
              env:
                  GITHUB_TOKEN: $__GHA_OPEN__ secrets.GITHUB_TOKEN __GHA_CLOSE__

    release:
        name: Release
        needs:
            - generate
            - lint
            - test
            - toolchain-security
            - security
        if: |
            github.event_name == 'push' &&
            needs.lint.result == 'success' &&
            needs.test.result == 'success' &&
            needs.toolchain-security.result == 'success' &&
            needs.security.result == 'success'
        permissions:
            contents: write
            issues: write
            pull-requests: write
        uses: ./.github/workflows/release.yml
        with:
            release_branch: $__GHA_OPEN__ github.ref_name __GHA_CLOSE__
            release_sha: $__GHA_OPEN__ github.sha __GHA_CLOSE__
        secrets: inherit

    promotion-pr:
        name: Create Promotion PR
        runs-on: ubuntu-latest
        needs:
            - release
        if: |
            github.event_name == 'push' &&
            (github.ref_name == 'develop' || github.ref_name == 'uat')
        permissions:
            contents: read
            issues: write
            pull-requests: write
        steps:
            - name: Create or update promotion PR
              env:
                  GH_TOKEN: $__GHA_OPEN__ github.token __GHA_CLOSE__
                  REPOSITORY: $__GHA_OPEN__ github.repository __GHA_CLOSE__
                  SOURCE_BRANCH: $__GHA_OPEN__ github.ref_name __GHA_CLOSE__
                  SOURCE_SHA: $__GHA_OPEN__ github.sha __GHA_CLOSE__
                  WORKFLOW_RUN_URL: $__GHA_OPEN__ github.server_url __GHA_CLOSE__/$__GHA_OPEN__ github.repository __GHA_CLOSE__/actions/runs/$__GHA_OPEN__ github.run_id __GHA_CLOSE__
              run: |
                  case "${SOURCE_BRANCH}" in
                      develop)
                          TARGET_BRANCH="uat"
                          ;;
                      uat)
                          TARGET_BRANCH="main"
                          ;;
                      *)
                          echo "No promotion target configured for ${SOURCE_BRANCH}."
                          exit 0
                          ;;
                  esac

                  compare_status="$(gh api \
                      "repos/${REPOSITORY}/compare/${TARGET_BRANCH}...${SOURCE_BRANCH}" \
                      --jq '.status')"
                  if [ "${compare_status}" = "identical" ] || [ "${compare_status}" = "behind" ]; then
                      echo "${SOURCE_BRANCH} has no new commits to promote to ${TARGET_BRANCH}."
                      exit 0
                  fi

                  title="Promote ${SOURCE_BRANCH} to ${TARGET_BRANCH}"
                  body="$(mktemp)"
                  cat > "${body}" <<EOF
                  Automated promotion PR from ${SOURCE_BRANCH} to ${TARGET_BRANCH}.

                  Source commit: ${SOURCE_SHA}
                  Successful workflow run: ${WORKFLOW_RUN_URL}
                  EOF

                  existing_prs="$(gh pr list \
                      --repo "${REPOSITORY}" \
                      --head "${SOURCE_BRANCH}" \
                      --base "${TARGET_BRANCH}" \
                      --state open \
                      --json number,createdAt \
                      --jq 'sort_by(.createdAt) | reverse | .[].number')"

                  existing_pr="$(printf '%s\n' "${existing_prs}" | sed -n '1p')"
                  duplicate_prs="$(printf '%s\n' "${existing_prs}" | sed -n '2,$p')"

                  if [ -n "${existing_pr}" ]; then
                      gh pr edit "${existing_pr}" \
                          --repo "${REPOSITORY}" \
                          --title "${title}" \
                          --body-file "${body}"
                      echo "Updated promotion PR #${existing_pr}."
                  else
                      created_pr_url="$(gh pr create \
                          --repo "${REPOSITORY}" \
                          --head "${SOURCE_BRANCH}" \
                          --base "${TARGET_BRANCH}" \
                          --title "${title}" \
                          --body-file "${body}")"
                      existing_pr="$(gh pr view "${created_pr_url}" \
                          --repo "${REPOSITORY}" \
                          --json number \
                          --jq '.number')"
                      echo "Created promotion PR #${existing_pr}."
                  fi

                  for duplicate_pr in ${duplicate_prs}; do
                      gh pr close "${duplicate_pr}" \
                          --repo "${REPOSITORY}" \
                          --comment "Closing duplicate promotion PR; #${existing_pr} now tracks ${SOURCE_BRANCH} to ${TARGET_BRANCH}."
                  done

    helm-image-update:
        name: Update Helm Values
        runs-on: ubuntu-latest
        needs:
            - release
        if: |
            github.event_name == 'push' &&
            needs.release.outputs.published == 'true'
        permissions:
            contents: read
        steps:
            # actions/create-github-app-token v2.1.4
            - uses: actions/create-github-app-token@67018539274d69449ef7c02e8e71183d1719ab42
              id: homelab-devops-token
              with:
                  app-id: $__GHA_OPEN__ secrets.BRANCH_SYNC_APP_ID __GHA_CLOSE__
                  private-key: $__GHA_OPEN__ secrets.BRANCH_SYNC_APP_PRIVATE_KEY __GHA_CLOSE__
                  owner: $__GHA_OPEN__ github.repository_owner __GHA_CLOSE__
                  repositories: homelab-devops

            # actions/checkout v6.0.2
            - uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd
              with:
                  repository: $__GHA_OPEN__ github.repository_owner __GHA_CLOSE__/homelab-devops
                  ref: main
                  path: homelab-devops
                  token: $__GHA_OPEN__ steps.homelab-devops-token.outputs.token __GHA_CLOSE__

            - name: Update environment values image
              env:
                  VALUES_FILE: homelab-devops/apps/kubernetes/app/$__GHA_OPEN__ github.ref_name __GHA_CLOSE__/{{ .Name }}/values.yml
                  IMAGE_REPOSITORY: zot.lan/$__GHA_OPEN__ github.repository __GHA_CLOSE__
                  RELEASE_TAG: $__GHA_OPEN__ needs.release.outputs.git_tag __GHA_CLOSE__
                  IMAGE_DIGEST: $__GHA_OPEN__ needs.release.outputs.image_digest __GHA_CLOSE__
              run: |
                  if [ ! -f "${VALUES_FILE}" ]; then
                      echo "${VALUES_FILE} does not exist; skipping Helm values update."
                      exit 0
                  fi

                  tmp_file="$(mktemp)"
                  awk \
                      -v image_repository="${IMAGE_REPOSITORY}" \
                      -v release_tag="${RELEASE_TAG}" \
                      -v image_digest="${IMAGE_DIGEST}" '
                      /^image:/ {
                          in_image = 1
                          print
                          next
                      }
                      in_image && $0 !~ /^    / {
                          in_image = 0
                      }
                      in_image && $1 == "repository:" {
                          print "    repository: " image_repository
                          next
                      }
                      in_image && $1 == "tag:" {
                          print "    tag: \"" release_tag "\""
                          next
                      }
                      in_image && $1 == "digest:" {
                          print "    digest: \"" image_digest "\""
                          next
                      }
                      { print }
                  ' "${VALUES_FILE}" > "${tmp_file}"
                  mv "${tmp_file}" "${VALUES_FILE}"

            - name: Commit values update
              env:
                  ENVIRONMENT: $__GHA_OPEN__ github.ref_name __GHA_CLOSE__
                  RELEASE_TAG: $__GHA_OPEN__ needs.release.outputs.git_tag __GHA_CLOSE__
              run: |
                  cd homelab-devops

                  if git diff --quiet; then
                      echo "${ENVIRONMENT} Helm values already point at ${RELEASE_TAG}."
                      exit 0
                  fi

                  git config user.name "homelab-image-update[bot]"
                  git config user.email "homelab-image-update[bot]@users.noreply.github.com"
                  git add "apps/kubernetes/app/${ENVIRONMENT}/{{ .Name }}/values.yml"
                  git commit -m "chore({{ .Name }}): update ${ENVIRONMENT} image to ${RELEASE_TAG}"

                  for attempt in 1 2 3; do
                      if git push origin HEAD:main; then
                          exit 0
                      fi

                      git fetch origin main
                      git rebase origin/main
                  done

                  echo "Unable to push Helm values update after retries."
                  exit 1

    sync-down:
        name: Sync Main To Prerelease Branches
        runs-on: ubuntu-latest
        needs:
            - release
        if: |
            github.event_name == 'push' &&
            github.ref_name == 'main'
        permissions:
            contents: read
        steps:
            # actions/create-github-app-token v2.1.4
            - uses: actions/create-github-app-token@67018539274d69449ef7c02e8e71183d1719ab42
              id: branch-sync-token
              with:
                  app-id: $__GHA_OPEN__ secrets.BRANCH_SYNC_APP_ID __GHA_CLOSE__
                  private-key: $__GHA_OPEN__ secrets.BRANCH_SYNC_APP_PRIVATE_KEY __GHA_CLOSE__

            # actions/checkout v6.0.2
            - uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd
              with:
                  fetch-depth: 0
                  token: $__GHA_OPEN__ steps.branch-sync-token.outputs.token __GHA_CLOSE__

            - name: Fast-forward prerelease branches
              env:
                  REPOSITORY: $__GHA_OPEN__ github.repository __GHA_CLOSE__
                  BRANCH_SYNC_TOKEN: $__GHA_OPEN__ steps.branch-sync-token.outputs.token __GHA_CLOSE__
              run: |
                  git config user.name "homelab-branch-sync[bot]"
                  git config user.email "homelab-branch-sync[bot]@users.noreply.github.com"
                  git remote set-url origin "https://x-access-token:${BRANCH_SYNC_TOKEN}@github.com/${REPOSITORY}.git"

                  git fetch origin main uat develop

                  git checkout -B uat origin/uat
                  git merge --ff-only origin/main
                  git push origin uat

                  git checkout -B develop origin/develop
                  git merge --ff-only uat
                  git push origin develop
	`,
	".github/workflows/release.yml": `name: Release

on:
    workflow_call:
        inputs:
            release_branch:
                required: true
                type: string
            release_sha:
                required: true
                type: string
        outputs:
            published:
                value: $__GHA_OPEN__ jobs.version.outputs.published __GHA_CLOSE__
            version:
                value: $__GHA_OPEN__ jobs.version.outputs.version __GHA_CLOSE__
            git_tag:
                value: $__GHA_OPEN__ jobs.version.outputs.git_tag __GHA_CLOSE__
            image_digest:
                value: $__GHA_OPEN__ jobs.publish.outputs.image_digest __GHA_CLOSE__
            arm64_digest:
                value: $__GHA_OPEN__ jobs.publish.outputs.arm64_digest __GHA_CLOSE__

concurrency:
    group: $__GHA_OPEN__ github.workflow __GHA_CLOSE__-publish
    cancel-in-progress: false

permissions:
    contents: read

env:
    REGISTRY: zot.kittiaccess.work
    RELEASE_BRANCH: $__GHA_OPEN__ inputs.release_branch __GHA_CLOSE__
    RELEASE_SHA: $__GHA_OPEN__ inputs.release_sha __GHA_CLOSE__

jobs:
    version:
        name: Publish Semantic Release
        runs-on: ubuntu-latest
        permissions:
            contents: write
            issues: write
            pull-requests: write
        outputs:
            published: $__GHA_OPEN__ steps.semantic.outputs.published __GHA_CLOSE__
            version: $__GHA_OPEN__ steps.semantic.outputs.version __GHA_CLOSE__
            git_tag: $__GHA_OPEN__ steps.semantic.outputs.git_tag __GHA_CLOSE__
            env_tag: $__GHA_OPEN__ steps.semantic.outputs.env_tag __GHA_CLOSE__
        steps:
            # actions/checkout v6.0.2
            - uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd
              with:
                  ref: $__GHA_OPEN__ env.RELEASE_BRANCH __GHA_CLOSE__
                  fetch-depth: 0

            - name: Pin release commit
              run: git reset --hard "${RELEASE_SHA}"

            - name: Install semantic-release
              run: |
                  npm install --no-save --no-package-lock \
                      semantic-release@24.2.7 \
                      @semantic-release/commit-analyzer@13.0.1 \
                      @semantic-release/release-notes-generator@14.0.3 \
                      @semantic-release/github@11.0.6 \
                      conventional-changelog-conventionalcommits@8.0.0

            - name: Publish semantic release
              id: semantic
              env:
                  GITHUB_TOKEN: $__GHA_OPEN__ secrets.GITHUB_TOKEN __GHA_CLOSE__
              run: |
                  echo "published=false" >> "${GITHUB_OUTPUT}"
                  echo "env_tag=${RELEASE_BRANCH}" >> "${GITHUB_OUTPUT}"

                  npx semantic-release --no-ci > "${RUNNER_TEMP}/semantic-release.log" 2>&1
                  cat "${RUNNER_TEMP}/semantic-release.log"

                  version="$(sed -nE \
                      -e 's/.*Published release ([0-9A-Za-z.+-]+).*/\1/p' \
                      -e 's/.*The next release version is ([0-9A-Za-z.+-]+).*/\1/p' \
                      "${RUNNER_TEMP}/semantic-release.log" | tail -n 1)"
                  if [ -z "${version}" ]; then
                      version="$(git tag --points-at HEAD --sort=-version:refname | sed -nE 's/^v([0-9A-Za-z.+-]+)$/\1/p' | head -n 1)"
                  fi
                  if [ -z "${version}" ]; then
                      echo "No semantic release will be published."
                      exit 0
                  fi

                  echo "published=true" >> "${GITHUB_OUTPUT}"
                  echo "version=${version}" >> "${GITHUB_OUTPUT}"
                  echo "git_tag=v${version}" >> "${GITHUB_OUTPUT}"

    build:
        name: Build and Scan Image ($__GHA_OPEN__ matrix.arch __GHA_CLOSE__)
        runs-on: $__GHA_OPEN__ matrix.runner __GHA_CLOSE__
        needs:
            - version
        if: needs.version.outputs.published == 'true'
        strategy:
            fail-fast: false
            matrix:
                include:
                    # Re-enable when x86 images are needed again.
                    # - arch: amd64
                    #   platform: linux/amd64
                    #   runner: ubuntu-latest
                    - arch: arm64
                      platform: linux/arm64
                      runner: ubuntu-24.04-arm
        env:
            IMAGE_REF: zot.kittiaccess.work/$__GHA_OPEN__ github.repository __GHA_CLOSE__
            STAGING_IMAGE_REF: zot.kittiaccess.work/$__GHA_OPEN__ github.repository __GHA_CLOSE__:staging-$__GHA_OPEN__ inputs.release_sha __GHA_CLOSE__-$__GHA_OPEN__ matrix.arch __GHA_CLOSE__
            ARCH_IMAGE_REF: zot.kittiaccess.work/$__GHA_OPEN__ github.repository __GHA_CLOSE__:$__GHA_OPEN__ needs.version.outputs.git_tag __GHA_CLOSE__-$__GHA_OPEN__ matrix.arch __GHA_CLOSE__
            SHA_ARCH_IMAGE_REF: zot.kittiaccess.work/$__GHA_OPEN__ github.repository __GHA_CLOSE__:$__GHA_OPEN__ inputs.release_sha __GHA_CLOSE__-$__GHA_OPEN__ matrix.arch __GHA_CLOSE__
        steps:
            # actions/checkout v6.0.2
            - uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd
              with:
                  ref: $__GHA_OPEN__ env.RELEASE_SHA __GHA_CLOSE__

            # docker/setup-buildx-action v3.11.1
            - uses: docker/setup-buildx-action@e468171a9de216ec08956ac3ada2f0791b6bd435

            # docker/login-action v4.1.0
            - name: Login to Zot registry
              uses: docker/login-action@4907a6ddec9925e35a0a9e82d7399ccc52663121
              with:
                  registry: $__GHA_OPEN__ env.REGISTRY __GHA_CLOSE__
                  username: $__GHA_OPEN__ secrets.ZOT_USERNAME __GHA_CLOSE__
                  password: $__GHA_OPEN__ secrets.ZOT_TOKEN __GHA_CLOSE__

            # docker/build-push-action v6.18.0
            - name: Build and push staging Docker image
              uses: docker/build-push-action@263435318d21b8e681c14492fe198d362a7d2c83
              with:
                  context: .
                  platforms: $__GHA_OPEN__ matrix.platform __GHA_CLOSE__
                  tags: $__GHA_OPEN__ env.STAGING_IMAGE_REF __GHA_CLOSE__
                  outputs: type=image,push=true,oci-mediatypes=true
                  provenance: true
                  sbom: true

            # aquasecurity/trivy-action v0.36.0
            - uses: aquasecurity/trivy-action@a9c7b0f06e461e9d4b4d1711f154ee024b8d7ab8
              with:
                  image-ref: $__GHA_OPEN__ env.STAGING_IMAGE_REF __GHA_CLOSE__
                  format: table
                  exit-code: 1
                  severity: CRITICAL,HIGH
                  ignore-unfixed: true
              env:
                  TRIVY_IMAGE_SRC: remote
                  TRIVY_PLATFORM: $__GHA_OPEN__ matrix.platform __GHA_CLOSE__
                  TRIVY_USERNAME: $__GHA_OPEN__ secrets.ZOT_USERNAME __GHA_CLOSE__
                  TRIVY_PASSWORD: $__GHA_OPEN__ secrets.ZOT_TOKEN __GHA_CLOSE__

            - name: Promote scanned image tags
              run: |
                  docker buildx imagetools create \
                      --tag "${ARCH_IMAGE_REF}" \
                      --tag "${SHA_ARCH_IMAGE_REF}" \
                      "${STAGING_IMAGE_REF}"

    publish:
        name: Publish Multi-Arch Image
        runs-on: ubuntu-latest
        needs:
            - version
            - build
        if: needs.version.outputs.published == 'true'
        outputs:
            image_digest: $__GHA_OPEN__ steps.image.outputs.digest __GHA_CLOSE__
            arm64_digest: $__GHA_OPEN__ steps.image.outputs.arm64_digest __GHA_CLOSE__
        env:
            IMAGE_REF: zot.kittiaccess.work/$__GHA_OPEN__ github.repository __GHA_CLOSE__
            RELEASE_TAG: $__GHA_OPEN__ needs.version.outputs.git_tag __GHA_CLOSE__
            ENV_TAG: $__GHA_OPEN__ needs.version.outputs.env_tag __GHA_CLOSE__
        steps:
            # docker/setup-buildx-action v3.11.1
            - uses: docker/setup-buildx-action@e468171a9de216ec08956ac3ada2f0791b6bd435

            # docker/login-action v4.1.0
            - name: Login to Zot registry
              uses: docker/login-action@4907a6ddec9925e35a0a9e82d7399ccc52663121
              with:
                  registry: $__GHA_OPEN__ env.REGISTRY __GHA_CLOSE__
                  username: $__GHA_OPEN__ secrets.ZOT_USERNAME __GHA_CLOSE__
                  password: $__GHA_OPEN__ secrets.ZOT_TOKEN __GHA_CLOSE__

            - name: Create multi-arch manifest
              run: |
                  tags=(
                      --tag "${IMAGE_REF}:${RELEASE_TAG}"
                      --tag "${IMAGE_REF}:${ENV_TAG}"
                      --tag "${IMAGE_REF}:${RELEASE_SHA}"
                  )
                  if [ "${RELEASE_BRANCH}" = "main" ]; then
                      tags+=(--tag "${IMAGE_REF}:latest")
                  fi

                  images=(
                      # "${IMAGE_REF}:${RELEASE_TAG}-amd64"
                      "${IMAGE_REF}:${RELEASE_TAG}-arm64"
                  )

                  docker buildx imagetools create \
                      "${tags[@]}" \
                      "${images[@]}"

            - name: Resolve image digest
              id: image
              run: |
                  digest="$(docker buildx imagetools inspect "${IMAGE_REF}:${RELEASE_TAG}" --format '__TPL_OPEN__.Manifest.Digest__TPL_CLOSE__')"
                  # amd64_digest="$(docker buildx imagetools inspect "${IMAGE_REF}:${RELEASE_TAG}-amd64" --format '__TPL_OPEN__.Manifest.Digest__TPL_CLOSE__')"
                  arm64_digest="$(docker buildx imagetools inspect "${IMAGE_REF}:${RELEASE_TAG}-arm64" --format '__TPL_OPEN__.Manifest.Digest__TPL_CLOSE__')"
                  echo "digest=${digest}" >> "${GITHUB_OUTPUT}"
                  # echo "amd64_digest=${amd64_digest}" >> "${GITHUB_OUTPUT}"
                  echo "arm64_digest=${arm64_digest}" >> "${GITHUB_OUTPUT}"

            # sigstore/cosign-installer v4.1.1
            - uses: sigstore/cosign-installer@cad07c2e89fa2edd6e2d7bab4c1aa38e53f76003

            - name: Sign Docker images
              env:
                  COSIGN_PRIVATE_KEY: $__GHA_OPEN__ secrets.COSIGN_PRIVATE_KEY __GHA_CLOSE__
              run: |
                  key_file="$(mktemp)"
                  trap 'rm -f "${key_file}"' EXIT
                  printf '%s' "${COSIGN_PRIVATE_KEY}" > "${key_file}"
                  chmod 600 "${key_file}"

                  cosign sign --yes \
                      --new-bundle-format=false \
                      --use-signing-config=false \
                      --key "${key_file}" \
                      "${IMAGE_REF}@$__GHA_OPEN__ steps.image.outputs.digest __GHA_CLOSE__"

                  # cosign sign --yes \
                  #     --new-bundle-format=false \
                  #     --use-signing-config=false \
                  #     --key "${key_file}" \
                  #     "${IMAGE_REF}@$__GHA_OPEN__ steps.image.outputs.amd64_digest __GHA_CLOSE__"

                  cosign sign --yes \
                      --new-bundle-format=false \
                      --use-signing-config=false \
                      --key "${key_file}" \
                      "${IMAGE_REF}@$__GHA_OPEN__ steps.image.outputs.arm64_digest __GHA_CLOSE__"
`,
	".releaserc.json": `{
    "branches": [
        {
            "name": "develop",
            "prerelease": "beta"
        },
        {
            "name": "uat",
            "prerelease": "rc"
        },
        "main"
    ],
    "tagFormat": "v${version}",
    "plugins": [
        [
            "@semantic-release/commit-analyzer",
            {
                "preset": "conventionalcommits",
                "parserOpts": {
                    "noteKeywords": [
                        "BREAKING CHANGE",
                        "BREAKING CHANGES",
                        "BREAKING"
                    ]
                },
                "releaseRules": [
                    {
                        "type": "build",
                        "release": false
                    },
                    {
                        "type": "ci",
                        "release": false
                    },
                    {
                        "type": "docs",
                        "release": false
                    },
                    {
                        "type": "style",
                        "release": false
                    },
                    {
                        "type": "test",
                        "release": false
                    },
                    {
                        "type": "chore",
                        "release": false
                    },
                    {
                        "type": "refactor",
                        "release": "patch"
                    },
                    {
                        "type": "perf",
                        "release": "patch"
                    }
                ]
            }
        ],
        [
            "@semantic-release/release-notes-generator",
            {
                "preset": "conventionalcommits",
                "parserOpts": {
                    "noteKeywords": [
                        "BREAKING CHANGE",
                        "BREAKING CHANGES",
                        "BREAKING"
                    ]
                }
            }
        ],
        "@semantic-release/github"
    ]
}
`,
	"Dockerfile": `FROM zot.kittiaccess.work/kitti12911/image-toolchain@sha256:47355f96a059465947c38aa956da1c4502c11d1e8f53eb2c8b3980ba58983d42 AS builder

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY buf.gen.yaml ./
COPY cmd ./cmd
COPY internal ./internal
COPY proto ./proto

RUN rm -rf gen/grpc \
	&& buf generate

ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
	go build -trimpath -ldflags="-s -w" -o /out/{{ .BinaryName }} ./cmd/server

FROM alpine:3.22@sha256:310c62b5e7ca5b08167e4384c68db0fd2905dd9c7493756d356e893909057601

RUN apk add --no-cache ca-certificates tzdata \
	&& addgroup -S app \
	&& adduser -S -G app app

WORKDIR /app

COPY --from=builder /out/{{ .BinaryName }} /app/{{ .BinaryName }}
COPY --chown=app:app config.example.yml /app/config.yml

USER app

EXPOSE {{ .GRPCPort }}

ENTRYPOINT ["/app/{{ .BinaryName }}"]
`,
	"Makefile": `# ____________________ Go Command ____________________
air:
	air

tidy:
	go mod tidy

run:
	go run ./cmd/server

lint: gen vet golangci-lint markdownlint lint-proto

vet: gen
	go vet ./...

golangci-lint: gen
	golangci-lint run --timeout=5m

markdownlint:
	markdownlint-cli2

lint-proto:
	buf lint

fmt:
	go fmt ./...

pretty:
	prettier --write "**/*.{md,markdown,yml,yaml,json,jsonc}"

format-proto:
	buf format -w

format: gen fmt pretty format-proto

test: gen
	env CGO_ENABLED=1 go test --race -v ./...

cov: gen
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

fix:
	go fix ./...

# ____________________ Generate Command ____________________
gen: gen-proto

gen-proto:
	rm -rf gen/grpc
	buf generate
`,
	"README.md": `# {{ .Name }}

Internal gRPC service bootstrap generated by service-gen.

## Local Development

Copy the example config when you want a local file:

` + "```sh" + `
cp config.example.yml config.yml
` + "```" + `

Run the service:

` + "```sh" + `
make tidy
make run
` + "```" + `

Run verification:

` + "```sh" + `
make gen
make format
make lint
make test
` + "```" + `

## Protobuf

This starter includes an embedded ` + "`proto/`" + ` directory so the generated
repository is ready to run the same generated-code and deploy pipeline used by
homelab gRPC services. Replace the starter proto with an external protobuf
source when the real service contract exists, then update ` + "`buf.gen.yaml`" + `
to point at that source.

## Configuration

The starter code reads configuration from environment variables first. The
` + "`config.example.yml`" + ` file documents the values expected by deployment
systems and humans.
`,
	"config.example.yml": `service:
    name: {{ .Name }}
    port: 50051
    shutdown_timeout: 10s

logging:
    level: warn
    add_source: false
    include_trace_id: true

tracing:
    enabled: true
    endpoint: alloy.lan:4317
    protocol: grpc
    insecure: false
    sample_ratio: 1.0

profiling:
    enabled: true
    server_address: https://pyroscope-ingest.lan
    namespace: {{ .Name }}
    basic_auth_user: ""
    basic_auth_password: ""
    tenant_id: ""

database:
    host: postgres.lan
    port: "5432"
    user: example_user
    password: example_password
    database: example
    run_migrations: false
    run_seeders: false
    pool:
        max_conns: 20
        min_conns: 5
        max_conn_life_time: 1h
        max_conn_idle_time: 6h
`,
	"internal/config/config.go": `package config

import (
	"time"

	"github.com/kitti12911/lib-monitor/profiling"
	"github.com/kitti12911/lib-monitor/tracing"
	liborm "github.com/kitti12911/lib-orm/v2"
	"github.com/kitti12911/lib-util/v3/logger"
)

// Config contains runtime settings for the service.
type Config struct {
	Service   Service          ` + "`mapstructure:\"service\" validate:\"required\"`" + `
	Logging   logger.Config    ` + "`mapstructure:\"logging\"`" + `
	Tracing   tracing.Config   ` + "`mapstructure:\"tracing\"`" + `
	Profiling profiling.Config ` + "`mapstructure:\"profiling\"`" + `
	Database  liborm.Config    ` + "`mapstructure:\"database\" validate:\"required\"`" + `
}

// Service contains service identity and listener settings.
type Service struct {
	Name            string        ` + "`mapstructure:\"name\"             env:\"SERVICE_NAME\"      validate:\"required\"`" + `
	Port            int           ` + "`mapstructure:\"port\"             env:\"PORT\"              validate:\"required,gte=1,lte=65535\"`" + `
	ShutdownTimeout time.Duration ` + "`mapstructure:\"shutdown_timeout\" env:\"SHUTDOWN_TIMEOUT\"`" + `
}
`,
	"internal/database/database.go": `package database

import (
	"context"

	"{{ .ModulePath }}/internal/config"

	orm "github.com/kitti12911/lib-orm/v2"
)

// New creates the service-owned database connection.
func New(ctx context.Context, cfg *config.Config) (*orm.DB, error) {
	db, err := orm.New(
		ctx,
		cfg.Database,
		orm.WithApplicationName(cfg.Service.Name),
		orm.WithModels(models()...),
		orm.WithTracing(cfg.Tracing.Enabled),
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
`,
	"internal/database/model.go": `package database

func models() []any {
	return nil
}
`,
}

var grpcTemplates = map[string]string{
	"go.mod": `module {{ .ModulePath }}

go 1.26.3

require (
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.3
	github.com/kitti12911/lib-monitor v1.8.0
	github.com/kitti12911/lib-orm/v2 v2.5.0
	github.com/kitti12911/lib-util/v3 v3.7.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.68.0
	google.golang.org/grpc v1.80.0
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/fsnotify/fsnotify v1.10.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.13 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.30.2 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grafana/pyroscope-go v1.2.8 // indirect
	github.com/grafana/pyroscope-go/godeltaprof v0.1.9 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.28.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/klauspost/compress v1.17.8 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/pelletier/go-toml/v2 v2.3.0 // indirect
	github.com/puzpuzpuz/xsync/v3 v3.5.1 // indirect
	github.com/sagikazarmark/locafero v0.12.0 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/spf13/viper v1.21.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc // indirect
	github.com/uptrace/bun v1.2.18 // indirect
	github.com/uptrace/bun/dbfixture v1.2.18 // indirect
	github.com/uptrace/bun/dialect/pgdialect v1.2.18 // indirect
	github.com/uptrace/bun/driver/pgdriver v1.2.18 // indirect
	github.com/uptrace/bun/extra/bunotel v1.2.18 // indirect
	github.com/uptrace/opentelemetry-go-extra/otelsql v0.3.2 // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.43.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.43.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.43.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.43.0 // indirect
	go.opentelemetry.io/otel/metric v1.43.0 // indirect
	go.opentelemetry.io/otel/sdk v1.43.0 // indirect
	go.opentelemetry.io/otel/trace v1.43.0 // indirect
	go.opentelemetry.io/proto/otlp v1.10.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/crypto v0.50.0 // indirect
	golang.org/x/net v0.52.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260401024825-9d38bb4040a9 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260406210006-6f92a3bedf2d // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	mellium.im/sasl v0.3.2 // indirect
)
`,
	"buf.gen.yaml": `version: v2
inputs:
    - directory: proto
managed:
    enabled: true
    override:
        - file_option: go_package_prefix
          value: {{ .ModulePath }}/gen/grpc
plugins:
    - local: protoc-gen-go
      out: gen/grpc
      opt: paths=source_relative
    - local: protoc-gen-go-grpc
      out: gen/grpc
      opt: paths=source_relative
`,
	"buf.yaml": `version: v2
modules:
    - path: proto
lint:
    use:
        - STANDARD
breaking:
    use:
        - FILE
`,
	"proto/{{ .ProtoPackagePath }}/v1/starter.proto": `syntax = "proto3";

package {{ .ProtoPackage }}.v1;

message PingRequest {
  string message = 1;
}

message PingResponse {
  string message = 1;
}

service StarterService {
  rpc Ping(PingRequest) returns (PingResponse);
}
`,
	"cmd/server/main.go": `package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"{{ .ModulePath }}/internal/config"
	"{{ .ModulePath }}/internal/database"
	"{{ .ModulePath }}/internal/server"

	"github.com/kitti12911/lib-monitor/profiling"
	"github.com/kitti12911/lib-monitor/tracing"
	libconfig "github.com/kitti12911/lib-util/v3/config"
	"github.com/kitti12911/lib-util/v3/logger"
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx := context.Background()

	cfg, err := libconfig.Load[config.Config]("config.yml")
	if err != nil {
		slog.ErrorContext(ctx, "failed to load config", "error", err)
		return 1
	}

	if cfg.Service.ShutdownTimeout == 0 {
		cfg.Service.ShutdownTimeout = 10 * time.Second
	}

	logger.NewFromConfig(cfg.Logging, cfg.Service.Name)

	profiler, err := profiling.NewFromConfig(cfg.Service.Name, cfg.Profiling)
	if err != nil {
		slog.ErrorContext(ctx, "failed to init profiling", "error", err)
		return 1
	}
	defer func() {
		if shutdownErr := profiling.Shutdown(profiler); shutdownErr != nil {
			slog.ErrorContext(ctx, "failed to stop profiling", "error", shutdownErr)
		}
	}()

	tp, err := tracing.NewFromConfig(ctx, cfg.Service.Name, cfg.Tracing)
	if err != nil {
		slog.ErrorContext(ctx, "failed to init tracing", "error", err)
		return 1
	}
	defer func() {
		if shutdownErr := tracing.Shutdown(ctx, tp); shutdownErr != nil {
			slog.ErrorContext(ctx, "failed to stop tracing", "error", shutdownErr)
		}
	}()

	db, err := database.New(ctx, cfg)
	if err != nil {
		slog.ErrorContext(ctx, "failed to init database", "error", err)
		return 1
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			slog.ErrorContext(ctx, "failed to close database", "error", closeErr)
		}
	}()

	srv, err := server.NewGRPCServer(ctx, cfg.Service.Port)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create gRPC server", "error", err)
		return 1
	}

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- srv.Start()
	}()

	slog.InfoContext(ctx, "gRPC server started", "port", cfg.Service.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
	case err := <-serverErr:
		slog.ErrorContext(ctx, "gRPC server error", "error", err)
		return 1
	}

	slog.InfoContext(ctx, "shutting down gRPC server")

	shutdownCtx, cancel := context.WithTimeout(ctx, cfg.Service.ShutdownTimeout)
	defer cancel()

	srv.Stop(shutdownCtx)

	slog.InfoContext(ctx, "server stopped")

	return 0
}
`,
	"internal/server/grpc.go": `package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	starterv1 "{{ .ModulePath }}/gen/grpc/{{ .ProtoPackagePath }}/v1"
	"{{ .ModulePath }}/internal/feature/starter"
	"{{ .ModulePath }}/internal/server/interceptor"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// GRPCServer owns the gRPC listener and health state.
type GRPCServer struct {
	server   *grpc.Server
	listener net.Listener
	health   *health.Server
}

// NewGRPCServer creates the gRPC server and base middleware.
func NewGRPCServer(ctx context.Context, port int) (*GRPCServer, error) {
	listenConfig := net.ListenConfig{}
	listener, err := listenConfig.Listen(ctx, "tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port %d: %w", port, err)
	}

	recoveryOpt := recovery.WithRecoveryHandler(func(p any) error {
		return status.Errorf(codes.Internal, "internal server error: %v", p)
	})

	srv := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler(
			otelgrpc.WithFilter(interceptor.TraceableRPC),
		)),
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recoveryOpt),
		),
		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(recoveryOpt),
		),
	)

	healthServer := health.NewServer()
	healthv1.RegisterHealthServer(srv, healthServer)
	healthServer.SetServingStatus("", healthv1.HealthCheckResponse_SERVING)
	starterv1.RegisterStarterServiceServer(srv, starter.NewHandler())

	reflection.Register(srv)

	return &GRPCServer{
		server:   srv,
		listener: listener,
		health:   healthServer,
	}, nil
}

// Start begins serving gRPC requests.
func (s *GRPCServer) Start() error {
	slog.Info("gRPC server listening", "addr", s.listener.Addr().String())
	return s.server.Serve(s.listener)
}

// Stop gracefully stops the gRPC server.
func (s *GRPCServer) Stop(ctx context.Context) {
	s.health.Shutdown()

	done := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		slog.WarnContext(ctx, "graceful shutdown timed out, forcing stop")
		s.server.Stop()
	}
}
`,
	"internal/server/interceptor/interceptor.go": `package interceptor

import (
	"strings"

	"google.golang.org/grpc/stats"
)

// TraceableRPC reports whether a gRPC method should emit tracing spans.
func TraceableRPC(info *stats.RPCTagInfo) bool {
	if info == nil {
		return true
	}

	return !IsHealthCheck(info.FullMethodName)
}

// IsHealthCheck reports whether the full gRPC method belongs to the health service.
func IsHealthCheck(fullMethod string) bool {
	return strings.HasPrefix(fullMethod, "/grpc.health.v1.Health/")
}
`,
	"internal/feature/starter/handler.go": `package starter

import (
	"context"

	starterv1 "{{ .ModulePath }}/gen/grpc/{{ .ProtoPackagePath }}/v1"
)

// Handler implements the embedded starter service.
type Handler struct {
	starterv1.UnimplementedStarterServiceServer
}

// NewHandler creates the starter gRPC handler.
func NewHandler() *Handler {
	return &Handler{}
}

// Ping echoes a message so the generated service has one callable RPC.
func (h *Handler) Ping(_ context.Context, req *starterv1.PingRequest) (*starterv1.PingResponse, error) {
	message := req.GetMessage()
	if message == "" {
		message = "pong"
	}

	return &starterv1.PingResponse{
		Message: message,
	}, nil
}
`,
}
