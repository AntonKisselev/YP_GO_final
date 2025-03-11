package task

import (
	"errors"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	Repeat string
}

func (t Task) NextDate(now time.Time, date string) (string, error) {
	re := regexp.MustCompile("^[0-9]{8}$")
	resBool := re.MatchString(date)
	if !resBool {
		return "", errors.New("invalid date")
	}
	startDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	repeat := t.Repeat
	if repeat == "" {
		return "", nil
	}

	/// d 7 format
	re = regexp.MustCompile(`^d ([0-9]+)$`)
	res := re.FindSubmatch([]byte(repeat))
	if res != nil {
		days, err := strconv.Atoi(string(res[1]))
		if err != nil {
			return "", err
		}
		if days > 400 {
			return "", errors.New("days > 400")
		}

		resultTime := startDate.AddDate(0, 0, days)
		for resultTime.Format("20060102") < now.Format("20060102") {
			resultTime = resultTime.AddDate(0, 0, days)
		}
		return resultTime.Format("20060102"), nil
	}
	/// y format
	re = regexp.MustCompile(`^y$`)
	resBool = re.MatchString(repeat)
	if resBool {
		resultTime := startDate.AddDate(1, 0, 0)
		for resultTime.Format("20060102") < now.Format("20060102") {
			resultTime = resultTime.AddDate(1, 0, 0)
		}
		return resultTime.Format("20060102"), nil
	}
	/// w 7 format
	re = regexp.MustCompile(`^w ([1-7,]+)$`)
	res = re.FindSubmatch([]byte(repeat))
	if res != nil {
		wdS := strings.Split(string(res[1]), ",")
		wdI := make([]int, len(wdS))
		for i, wd := range wdS {
			num, err := strconv.Atoi(wd)
			if err != nil {
				return "", err
			}
			if num > 8 || num < 1 {
				return "", errors.New("invalid week day number")
			}
			wdI[i] = num
		}

		resultTime := startDate.AddDate(0, 0, 1)
		for {
			if resultTime.Format("20060102") > now.Format("20060102") {
				rwd := int(resultTime.Weekday())
				if rwd == 0 {
					rwd = 7
				}
				if slices.Contains(wdI, rwd) {
					return resultTime.Format("20060102"), nil
				}
			}
			resultTime = resultTime.AddDate(0, 0, 1)
		}
	}
	/// m 1,-1 2,8 format
	re = regexp.MustCompile(`^m ([0-9,\- ]+)$`)
	res = re.FindSubmatch([]byte(repeat))
	if res != nil {
		/// slice num months
		var mNums []int
		/// slice num days month
		var mDays []int

		tmp := strings.Split(string(res[1]), " ")
		if len(tmp) > 2 {
			return "", errors.New("invalid months format")
		}
		if len(tmp) == 1 {
			mNums = make([]int, 0)
		}
		if len(tmp) == 2 {
			nNumsS := strings.Split(tmp[1], ",")
			mNums = make([]int, len(nNumsS))
			for i, n := range nNumsS {
				num, err := strconv.Atoi(n)
				if err != nil {
					return "", err
				}
				if num < 1 || num > 12 {
					return "", errors.New("invalid months number")
				}
				mNums[i] = num
			}
		}
		mDaysS := strings.Split(tmp[0], ",")
		mDays = make([]int, len(mDaysS))
		for i, m := range mDaysS {
			num, err := strconv.Atoi(m)
			if err != nil {
				return "", err
			}
			if num == -1 || num == -2 || (num >= 1 && num <= 31) {
				mDays[i] = num
			} else {
				return "", errors.New("invalid number day of month")
			}
		}

		resultTime := startDate.AddDate(0, 0, 1)
		var numMonthChecked bool
		var numDayMonthChecked bool
		for {
			if resultTime.Format("20060102") > now.Format("20060102") {
				numMonthChecked = false
				numDayMonthChecked = false

				if len(mNums) == 0 {
					numMonthChecked = true
				} else {
					if slices.Contains(mNums, int(resultTime.Month())) {
						numMonthChecked = true
					}
				}

				t := time.Date(resultTime.Year(), resultTime.Month(), 32, 0, 0, 0, 0, time.UTC)
				daysInMonth := 32 - t.Day()

				if slices.Contains(mDays, -1) && resultTime.Day() == daysInMonth {
					numDayMonthChecked = true
				}
				if slices.Contains(mDays, -2) && resultTime.Day() == daysInMonth-1 {
					numDayMonthChecked = true
				}
				if slices.Contains(mDays, resultTime.Day()) {
					numDayMonthChecked = true
				}

				if numMonthChecked && numDayMonthChecked {
					return resultTime.Format("20060102"), nil
				}

			}
			resultTime = resultTime.AddDate(0, 0, 1)
		}

	}

	return "", errors.New("unsupported repeat")
}
