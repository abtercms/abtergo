package page_test

import (
	"testing"

	"github.com/abtergo/abtergo/pkg/page"
)

func TestUpdater_Transition(t *testing.T) {
	updater, _ := setupTestUpdater()

	type args struct {
		status  page.Status
		trigger page.Trigger
	}
	tests := []struct {
		name    string
		args    args
		want    page.Status
		wantErr bool
	}{
		{
			name: "draft -> active: OK",
			args: args{
				status:  page.Draft,
				trigger: page.Activate,
			},
			want:    page.Active,
			wantErr: false,
		},
		{
			name: "draft -> inactive: OK",
			args: args{
				status:  page.Draft,
				trigger: page.Inactivate,
			},
			want:    page.Inactive,
			wantErr: false,
		},
		{
			name: "active -> inactive: OK",
			args: args{
				status:  page.Active,
				trigger: page.Inactivate,
			},
			want:    page.Inactive,
			wantErr: false,
		},
		{
			name: "inactive -> active: OK",
			args: args{
				status:  page.Inactive,
				trigger: page.Activate,
			},
			want:    page.Active,
			wantErr: false,
		},
		{
			name: "inactive -> inactive: NOT OK",
			args: args{
				status:  page.Inactive,
				trigger: page.Inactivate,
			},
			want:    page.Inactive,
			wantErr: true,
		},
		{
			name: "active -> active: NOT OK",
			args: args{
				status:  page.Active,
				trigger: page.Activate,
			},
			want:    page.Active,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := updater.Transition(tt.args.status, tt.args.trigger)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Transition() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type testUpdaterDeps struct{}

func (tud *testUpdaterDeps) AssertExpectations(t *testing.T) {
}

func setupTestUpdater() (page.Updater, testUpdaterDeps) {
	testUpdater := page.NewUpdater()

	return testUpdater, testUpdaterDeps{}
}
