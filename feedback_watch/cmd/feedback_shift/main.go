package main

import (
	"context"
	"fmt"
	"log"
	"time"

	at "github.com/bearyinnovative/radagast/airtable"
	"github.com/bearyinnovative/radagast/bearychat"
	"github.com/bearyinnovative/radagast/config"
	"github.com/fabioberger/airtable-go"
)

func main() {
	ctx := context.Background()
	ctx = config.MustMakeContext(ctx, "./radagast.toml")
	ctx = bearychat.MustMakeContext(ctx)
	ctx = at.MustMakeContext(ctx)

	if err := ExecuteOnce(ctx); err != nil {
		log.Fatalf("execute task failed: %+v", err)
	}
}

type ShiftRecord struct {
	Fields struct {
		Date      string   `json:"值日日期"`
		InCharges []string `json:"值日同学"`
	} `json:"fields"`
}

type Shift struct {
	Date      time.Time
	InCharges []Watchman
}

func (s Shift) GetInCharges() string {
	inCharges := ""
	for _, inCharge := range s.InCharges {
		inCharges = fmt.Sprintf("%s%s ", inCharges, inCharge.Name)
	}

	return fmt.Sprintf("`%s` %s", s.Date.Format("2006-01-02"), inCharges)
}

type WatchmenRecord struct {
	Id     string `json:"id"`
	Fields struct {
		Name string `json:"Name"`
	} `json:"fields"`
}

type Watchman struct {
	Id   string
	Name string
}

func ExecuteOnce(ctx context.Context) error {
	taskConfig := config.FromContext(ctx).Get("feedback_watch.shift").Config()

	airtableClient := at.ClientFromContext(ctx)

	watchmen, err := listWatchmen(ctx, taskConfig, airtableClient)
	if err != nil {
		return err
	}

	shifts, err := listShifts(ctx, taskConfig, airtableClient, watchmen)
	if err != nil {
		return err
	}

	todayShift, forcastShifts := planShifts(time.Now(), shifts)

	if todayShift == nil {
		return notifyMisconfig(ctx, taskConfig)
	} else {
		return sendPlan(ctx, taskConfig, todayShift, forcastShifts)
	}
}

func notifyMisconfig(ctx context.Context, config config.Config) error {
	return bearychat.SendToVchannel(
		ctx,
		bearychat.RTMClientFromContext(ctx),
		bearychat.RTMMessage{
			Text:       "好像今天没有人值班？",
			VchannelId: config.Get("misconfig-vchannel-id").String(),
		},
	)
}

func sendPlan(ctx context.Context, config config.Config, today *Shift, forcast []*Shift) error {
	forcastInCharges := ""
	if len(forcast) > 0 {
		for i, f := range forcast {
			if i > 2 {
				break
			}
			forcastInCharges = fmt.Sprintf(
				"%s%s | ",
				forcastInCharges,
				f.GetInCharges(),
			)
		}
	} else {
		forcastInCharges = "好像还没有排期哦"
	}

	return bearychat.SendToVchannel(
		ctx,
		bearychat.RTMClientFromContext(ctx),
		bearychat.RTMMessage{
			Text: fmt.Sprintf(
				config.Get("template").String(),
				today.GetInCharges(),
				forcastInCharges,
			),
			VchannelId: config.Get("feedback-vchannel-id").String(),
		},
	)
}

func listWatchmen(ctx context.Context, config config.Config, client *airtable.Client) (map[string]Watchman, error) {
	table := config.Get("incharge-table").String()

	watchmenRecords := []WatchmenRecord{}
	if err := client.ListRecords(
		table,
		&watchmenRecords,
		airtable.ListParameters{MaxRecords: 100},
	); err != nil {
		return nil, err
	}

	watchmen := make(map[string]Watchman)
	for _, watchmanRecord := range watchmenRecords {
		watchmen[watchmanRecord.Id] = Watchman{
			Id:   watchmanRecord.Id,
			Name: watchmanRecord.Fields.Name,
		}
	}

	return watchmen, nil
}

func listShifts(ctx context.Context, config config.Config, client *airtable.Client, watchmen map[string]Watchman) ([]*Shift, error) {
	table := config.Get("shift-table").String()
	view := config.Get("shift-table-view").String()

	shiftRecords := []ShiftRecord{}
	if err := client.ListRecords(
		table,
		&shiftRecords,
		airtable.ListParameters{View: view},
	); err != nil {
		return nil, err
	}

	var shifts []*Shift
	for _, shiftRecord := range shiftRecords {
		date, _ := time.Parse("2006-01-02", shiftRecord.Fields.Date)

		var inCharges []Watchman
		for _, id := range shiftRecord.Fields.InCharges {
			inCharges = append(inCharges, watchmen[id])
		}

		shifts = append(
			shifts,
			&Shift{
				Date:      date,
				InCharges: inCharges,
			},
		)
	}

	return shifts, nil
}

func dateBefore(a, b time.Time) bool {
	if a.Year() < b.Year() {
		return true
	}
	if a.Year() > b.Year() {
		return false
	}

	if a.Month() < b.Month() {
		return true
	}
	if a.Month() > b.Month() {
		return false
	}

	if a.Day() < b.Day() {
		return true
	}
	if a.Day() > b.Day() {
		return false
	}

	return false
}

func dateAfter(a, b time.Time) bool {
	if a.Year() > b.Year() {
		return true
	}
	if a.Year() < b.Year() {
		return false
	}

	if a.Month() > b.Month() {
		return true
	}
	if a.Month() < b.Month() {
		return false
	}

	if a.Day() > b.Day() {
		return true
	}
	if a.Day() < b.Day() {
		return false
	}

	return false
}

func planShifts(today time.Time, shifts []*Shift) (*Shift, []*Shift) {
	var (
		todayShift    *Shift
		forcastShifts []*Shift
	)

	for _, shift := range shifts {
		if dateBefore(shift.Date, today) {
			continue
		}

		if dateAfter(shift.Date, today) {
			forcastShifts = append(forcastShifts, shift)
			continue
		}

		todayShift = shift
	}

	return todayShift, forcastShifts
}
