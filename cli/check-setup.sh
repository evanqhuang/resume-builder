#!/bin/bash

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "Resume CLI - Setup Checker"
echo "=========================="
echo ""

# Check Go installation
echo -n "Checking Go installation... "
if command -v go &> /dev/null; then
    VERSION=$(go version | awk '{print $3}')
    echo -e "${GREEN}✓${NC} Found: $VERSION"
else
    echo -e "${RED}✗${NC} Go not found"
    echo "  Install from: https://go.dev/dl/"
    exit 1
fi

# Check pdflatex
echo -n "Checking pdflatex installation... "
if command -v pdflatex &> /dev/null; then
    echo -e "${GREEN}✓${NC} Found"
else
    echo -e "${YELLOW}⚠${NC} pdflatex not found (required for PDF generation)"
    echo "  macOS: brew install --cask mactex-no-gui"
    echo "  Linux: sudo apt-get install texlive-latex-base texlive-latex-extra"
fi

# Check resume.yaml
echo -n "Checking resume.yaml... "
RESUME_PATH="../resume.yaml"
if [ -f "$RESUME_PATH" ]; then
    echo -e "${GREEN}✓${NC} Found at $RESUME_PATH"
else
    echo -e "${RED}✗${NC} Not found at $RESUME_PATH"
    echo "  Update the path in your commands with -r flag"
fi

# Check .env file
echo -n "Checking .env configuration... "
if [ -f ".env" ]; then
    if grep -q "OPENROUTER_API_KEY" .env; then
        # Check if the value is not the placeholder
        if grep "OPENROUTER_API_KEY=your-api-key-here" .env &> /dev/null; then
            echo -e "${YELLOW}⚠${NC} .env exists but API key is placeholder"
            echo "  Edit .env and add your real OpenRouter API key"
        else
            echo -e "${GREEN}✓${NC} .env configured"
        fi
    else
        echo -e "${YELLOW}⚠${NC} .env missing OPENROUTER_API_KEY"
    fi
else
    echo -e "${YELLOW}⚠${NC} No .env file (required for 'match' command)"
    echo "  Copy .env.example to .env and add your API key"
fi

# Check if binary exists
echo -n "Checking binary... "
if [ -f "./resume-cli" ]; then
    echo -e "${GREEN}✓${NC} Found"
else
    echo -e "${YELLOW}⚠${NC} Binary not built yet"
    echo "  Run: make build"
fi

# Check dependencies
echo -n "Checking Go dependencies... "
if [ -f "go.sum" ]; then
    echo -e "${GREEN}✓${NC} Dependencies downloaded"
else
    echo -e "${YELLOW}⚠${NC} Dependencies not downloaded"
    echo "  Run: go mod download"
fi

echo ""
echo "Setup Summary"
echo "============="

# Count issues
ISSUES=0
if ! command -v go &> /dev/null; then ((ISSUES++)); fi
if ! command -v pdflatex &> /dev/null; then ((ISSUES++)); fi
if [ ! -f "$RESUME_PATH" ]; then ((ISSUES++)); fi
if [ ! -f ".env" ]; then ((ISSUES++)); fi
if [ ! -f "./resume-cli" ]; then ((ISSUES++)); fi

if [ $ISSUES -eq 0 ]; then
    echo -e "${GREEN}✓ All checks passed! You're ready to go.${NC}"
    echo ""
    echo "Try these commands:"
    echo "  ./resume-cli --help"
    echo "  ./resume-cli list"
    echo "  ./resume-cli generate -o test.pdf"
else
    echo -e "${YELLOW}⚠ Found $ISSUES issue(s). See above for details.${NC}"
fi

echo ""
